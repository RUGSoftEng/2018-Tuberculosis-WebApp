from django.conf.urls import url
from api import views

app_name = "api"


urlpatterns = [
    url(r'^$', views.patient_list, name="list"),
    url(r'^(?P<pk>\d+)/$', views.patient_details, name="patient_details"),
    url(r'^add/$', views.add_patient, name="add_patient"),
    url(r'^treatment/$', views.treatment, name="treatment"),
    url(r'^test/$', views.test, name="test"),
    url(r'^test2/$', views.test2, name="test2"),
    url(r'^login/$', views.login, name="login"),
    url(r'^register/$', views.register, name="register"),
    url(r'^logout/$', views.user_logout, name="logout"),
]