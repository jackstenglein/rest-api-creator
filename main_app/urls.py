from django.urls import path

from . import views


app_name = 'main_app'
urlpatterns = [
    path('objects', views.get_objects, name='get_objects'),
    path('objects/create', views.create_object, name='create_object'),
]
