from django.conf.urls import url
from api import views

app_name = "api"

urlpatterns = [
    url(r'^$', views.patient_list, name="list"),
    url(r'add$', views.add_patient, name="add_patient"),
    # url(r'^$', PatientList.as_view(), name="list"),
]