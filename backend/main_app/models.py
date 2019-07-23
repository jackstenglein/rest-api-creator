from django.db import models
# from django.contrib.auth.models import User
from django.conf import settings
from rest_framework import serializers
from enum import Enum, auto


class Project(models.Model):
    # Constants for Language choices
    SAILS_JS = 1
    LANGUAGE_CHOICES = (
        (SAILS_JS, 'Sails.js'),
    )

    name = models.CharField(max_length=30, blank=False)
    owner = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, blank=False)
    language = models.SmallIntegerField(choices=LANGUAGE_CHOICES, blank=False)

    def __str__(self):
        return str(ProjectSerializer(self).data)

    class Meta:
        db_table = 'project'

class ProjectSerializer(serializers.ModelSerializer):
    class Meta:
        model = Project
        fields = "__all__"


class Object(models.Model):
    project = models.ForeignKey(Project, on_delete=models.CASCADE, blank=False)
    owner = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, blank=False)
    name = models.CharField(max_length=30, blank=False)

    def __str__(self):
        return str(ObjectSerializer(self).data)

    class Meta:
        db_table = 'object'
        unique_together = (('name', 'project'),)

class ObjectSerializer(serializers.ModelSerializer):
    class Meta:
        model = Object
        fields = "__all__"


class Attribute(models.Model):
    STRING = 0
    INTEGER = STRING + 1
    FLOAT = INTEGER + 1
    DATETIME = FLOAT + 1
    BOOLEAN = DATETIME + 1
    LIST = BOOLEAN + 1
    DICT = LIST + 1
    DATA_TYPE_CHOICES = (
        (STRING, 'Text'),
        (INTEGER, 'Integer'),
        (FLOAT, 'Decimal'),
        (DATETIME, 'Date/Time'),
        (BOOLEAN, 'True/False'),
        (LIST, 'List'),
        (DICT, 'Dictionary')
    )
    DATA_TYPE_TO_SAILS_TYPE = (
        "string",
        "number",
        "number",
        "number",
        "boolean",
        "json",
        "json"
    )
    object = models.ForeignKey(Object, on_delete=models.CASCADE, blank=True, null=True)
    #action = models.ForeignKey(Action, on_delete=models.CASCADE, blank=True, null=True)
    name = models.CharField(max_length=30, blank=False)
    required = models.BooleanField(default=False)
    type = models.SmallIntegerField(choices=DATA_TYPE_CHOICES, blank=False)
    description = models.CharField(max_length=150, blank=True)

    def __str__(self):
        return str(AttributeSerializer(self).data)

    class Meta:
        db_table = 'attribute'
        unique_together = (('name', 'object'),)

class AttributeSerializer(serializers.ModelSerializer):
    class Meta:
        model = Attribute
        fields = ("id", "name", "required", "type")


class HttpMethods(Enum):
    POST = auto()
    GET = auto()
    PUT = auto()
    DELETE = auto()


class Action(models.Model):
    friendly_name = models.CharField(max_length=30, blank=False)
    kebab_case_name = models.CharField(max_length=30, blank=False)
    description = models.CharField(max_length=250, blank=True)
    object = models.ForeignKey(Object, on_delete=models.CASCADE, blank=False)
    route = models.CharField(max_length=100, blank=False)
    METHOD_TYPE_CHOICES = (
        (HttpMethods.POST, 'Create'),
        (HttpMethods.GET, 'Read'),
        (HttpMethods.PUT, 'Update'),
        (HttpMethods.DELETE, 'Delete')
    )
    method = models.SmallIntegerField(choices=METHOD_TYPE_CHOICES, blank=False)

    class Meta:
        db_table = 'action'
        unique_together = (('kebab_case_name', 'object'))


class ActionSerializer(serializers.ModelSerializer):
    class Meta:
        model = Action
        fields = "__all__"


class ReturnValueType(Enum):
    STATIC_STRING = auto()
    STATIC_NUMBER = auto()
    DATABASE_QUERY = auto()
    OBJECT_ATTRIBUTE = auto()


class SearchType(Enum):
    EQUAL = auto()
    GREATER_THAN = auto()
    GREATER_THAN_OR_EQUAL = auto()
    LESS_THAN = auto()
    LESS_THAN_OR_EQUAL = auto()
    NOT_EQUAL = auto()


class ReturnValue(models.Model):
    # STATIC_STRING = 0
    # STATIC_NUMBER = STATIC_STRING + 1
    # DATABASE_QUERY = STATIC_NUMBER + 1
    # OBJECT_ATTRIBUTE = DATABASE_QUERY + 1
    RETURN_TYPE_CHOICES = (
        (ReturnValueType.STATIC_STRING, "Static Text"),
        (ReturnValueType.STATIC_NUMBER, "Static Number"),
        (ReturnValueType.DATABASE_QUERY, "Result of Database Query"),
        (ReturnValueType.OBJECT_ATTRIBUTE, "Object Attribute")
    )
    action = models.ForeignKey(Action, on_delete=models.CASCADE, blank=False)
    description = models.CharField(max_length=150, blank=True)
    key = models.CharField(max_length=30, blank=False)
    type = models.SmallIntegerField(choices=RETURN_TYPE_CHOICES, blank=False)

    # Relevant for STATIC_STRING type
    static_string = models.CharField(max_length=100, blank=True)

    # Relevant for STATIC_NUMBER type
    static_number = models.FloatField(blank=True, null=True)

    # Relevant for DATABASE_QUERY/OBJECT_ATTRIBUTE type
    query_limit = models.PositiveSmallIntegerField(blank=True, null=True)
    skip = models.PositiveSmallIntegerField(blank=True, null=True)
    sort_field = models.ForeignKey(Attribute, on_delete=models.PROTECT, blank=True, null=True, related_name="+")
    ASC = "ASC"
    DESC = "DESC"
    SORT_DIR_CHOICES = (
        (ASC, "Ascending"),
        (DESC, "Descending")
    )
    sort_direction = models.CharField(max_length=4, choices=SORT_DIR_CHOICES, blank=True)
    search_attribute = models.ForeignKey(Attribute, on_delete=models.PROTECT, blank=True, null=True, related_name="+")
    SEARCH_TYPE_CHOICES = (
        (SearchType.EQUAL, "Equal To"),
        (SearchType.GREATER_THAN, "Greater Than"),
        (SearchType.GREATER_THAN_OR_EQUAL, "Greater Than or Equal To"),
        (SearchType.LESS_THAN, "Less Than"),
        (SearchType.LESS_THAN_OR_EQUAL, "Less Than or Equal To"),
        (SearchType.NOT_EQUAL, "Not Equal")
    )
    search_type = models.SmallIntegerField(choices=SEARCH_TYPE_CHOICES, blank=True, null=True)
    search_parameter = models.CharField(max_length=100, blank=True)

