from django.shortcuts import render
from django.views import generic
from django.views.generic.base import View
from django.http import HttpResponse

from .forms import PatientForm

def patient_list(request):
    return render(request, "base.html")


def add_patient(request):

    if request.method == 'GET':
        form = PatientForm(request.GET)
        if form.is_valid():
            name = form.cleaned_data['name']
            username = form.cleaned_data['username']
            password = form.cleaned_data['password']
            api_token = form.cleaned_data['api_token']
            print(name)
            print(username)
            print(password)
            print(api_token)
            return render(request, "base.html")

    else:
        print("wrong")
        form = PatientForm()

    return render(request, "addpatient.html", {'form': form})
