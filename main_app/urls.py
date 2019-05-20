from django.urls import path

from . import project_views, user_views, object_views, action_views


app_name = 'main_app'
urlpatterns = [
    path('signup', user_views.signup, name='signup'),
    path('login', user_views.app_login, name='login'),
    path('projects', project_views.projects, name='projects'),
    path('projects/<int:project_id>/download', project_views.download_project, name='download_project'),
    path('projects/<int:project>/objects', object_views.objects, name='many_objects'),
    path('projects/<int:project>/objects/<int:object>', object_views.objects, name='single_object'),
    path('projects/<int:project>/objects/<int:object>/actions', action_views.actions, name='many_actions')
]
