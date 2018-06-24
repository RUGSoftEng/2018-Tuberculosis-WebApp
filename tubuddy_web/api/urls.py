from django.conf.urls import url
from api import views

app_name = "api"


# urlpatterns = [
#     url(r'^$', views.patient_list, name="list"),
#     url(r'^(?P<pk>\d+)/$', views.patient_details, name="patient_details"),
#     url(r'^add/$', views.add_patient, name="add_patient"),
#     url(r'^treatment/$', views.treatment, name="treatment"),
#     url(r'^login/$', views.login, name="login"),
#     url(r'^register/$', views.register, name="register"),
#     url(r'^logout/$', views.user_logout, name="logout"),
#     url(r'^scheduled_dosage/$', views.scheduled_dosage, name="scheduled_dosage"),
#     url(r'^create_dosage/$', views.create_dosage, name="create_dosage")
# ]

# urlpatterns = [
#     url(r'^$', views.login, name="login"),
#     url(r'^list/$', views.patient_list, name="list"),
#     url(r'^list/(?P<pk>\d+)/$', views.patient_details, name="patient_details"),
#     url(r'^add/$', views.add_patient, name="add_patient"),
#     url(r'^treatment/$', views.treatment, name="treatment"),
#     url(r'^register/$', views.register, name="register"),
#     url(r'^logout/$', views.user_logout, name="logout"),
#     url(r'^scheduled_dosage/$', views.scheduled_dosage, name="scheduled_dosage"),
#     url(r'^create_dosage/$', views.create_dosage, name="create_dosage")
# ]


urlpatterns = [
    url(r'^$', views.login, name="login"),
    url(r'^list/$', views.patient_list, name="list"),
    url(r'^list/(?P<pk>\d+)/$', views.patient_details, name="patient_details"),
    url(r'^add/$', views.add_patient, name="add_patient"),
    url(r'^treatment/(?P<pk>\d+)/$', views.treatment, name="treatment"),
    url(r'^register/$', views.register, name="register"),
    url(r'^logout/$', views.user_logout, name="logout"),
    url(r'^create_dosage/(?P<pk>\d+)/$', views.create_dosage, name="create_dosage")
]