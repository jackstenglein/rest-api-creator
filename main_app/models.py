from django.db import models


class User(models.Model):
    username = models.CharField(max_length=30, unique=True, blank=False)
    password = models.CharField(max_length=64, blank=False)

    def __str__(self):
        return self.username

    class Meta:
        db_table = 'user'


class Project(models.Model):
    # Constants for Language choices
    SAILS_JS = 1
    LANGUAGE_CHOICES = (
        (SAILS_JS, 'Sails.js'),
    )

    name = models.CharField(max_length=30, blank=False)
    owner = models.ForeignKey(User, on_delete=models.CASCADE, db_column='owner_id', blank=False)
    language = models.SmallIntegerField(choices=LANGUAGE_CHOICES, blank=False)

    def __str__(self):
        return self.name

    class Meta:
        db_table = 'project'


class Object(models.Model):
    project = models.ForeignKey(Project, on_delete=models.CASCADE, db_column='project_id', blank=False)
    name = models.CharField(max_length=30, blank=False)

    def __str__(self):
        return self.name

    class Meta:
        db_table = 'object'
        unique_together = (('name', 'project'),)


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
    object = models.ForeignKey(Object, on_delete=models.CASCADE, db_column='object_id', blank=False)
    name = models.CharField(max_length=30, blank=False)
    required = models.BooleanField(default=False)
    type = models.SmallIntegerField(choices=DATA_TYPE_CHOICES, blank=False)

    def __str__(self):
        return self.name

    class Meta:
        db_table = 'attribute'
        unique_together = (('name', 'object'),)
