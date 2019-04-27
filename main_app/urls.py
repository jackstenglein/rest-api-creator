from django.urls import path

from . import project_views, user_views, object_views


app_name = 'main_app'
urlpatterns = [
    path('signup', user_views.signup, name='signup'),
    path('projects', project_views.projects, name='projects'),
    path('projects/<int:project>/objects', object_views.objects, name='objects'),
    path('projects/<int:project>/objects/<int:object>', object_views.objects, name='objects')
]
