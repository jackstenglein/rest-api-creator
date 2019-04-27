from django.db import models
# from django.contrib.auth.models import User
from django.conf import settings
from rest_framework import serializers


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
        return self.name

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
    STRING = 1
    INTEGER = 2
    FLOAT = 3
    DATETIME = 4
    BOOLEAN = 5
    JSON = 6
    DATA_TYPE_CHOICES = (
        (STRING, 'Text'),
        (INTEGER, 'Integer'),
        (FLOAT, 'Decimal'),
        (DATETIME, 'Date/Time'),
        (BOOLEAN, 'True/False'),
        (JSON, 'JSON'),
    )
    object = models.ForeignKey(Object, on_delete=models.CASCADE, blank=False)
    name = models.CharField(max_length=30, blank=False)
    required = models.BooleanField(default=False)
    type = models.SmallIntegerField(choices=DATA_TYPE_CHOICES, blank=False)

    def __str__(self):
        return self.name

    class Meta:
        db_table = 'attribute'
        unique_together = (('name', 'object'),)


class AttributeSerializer(serializers.ModelSerializer):
    class Meta:
        model = Attribute
        fields = ("id", "name", "required", "type")
