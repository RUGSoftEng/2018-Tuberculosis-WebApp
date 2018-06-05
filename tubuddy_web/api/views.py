# from django.contrib.auth.models import User
from django.contrib.auth import authenticate, login as auth_login, logout
from django.shortcuts import render
from django.views import generic
from django.views.generic.base import View
from django.http import HttpResponseRedirect, HttpResponse
from django.shortcuts import redirect
from django.urls import reverse
from .models import UserProfile, User
from django.contrib.auth.decorators import login_required, user_passes_test

from .forms import *

import requests

print ("++++++++++++++++++++++++++")
User.objects.all().delete()
UserProfile.objects.all().delete()
print ("+++ Cleared all entries ++")
print ("++++++++++++++++++++++++++")


# create_physician_url = "http://192.168.50.4:2002/api/accounts/physicians"
# data = {
#      "username": "House",
#      "name": "Gregory",
#      "password": "pass",
#      "api_token": "7777777",
#      "email": "example@gmail.com",
#      "creation_token": "HOUSE"
#     }
# create_physician_req = requests.request('PUT', url=create_physician_url, json=data)
#
# login_url = "http://192.168.50.4:2002/api/accounts/login"
# data_doctor = {
#     "username": "House",
#     "password": "pass",
# }
# login_doctor_req = requests.request('POST', login_url, json=data)
# print (login_doctor_req.text)





# add_patient_url = "http://192.168.50.4:2002/api/accounts/patients?token=HOUSE"
# data1 = {
#      "username": "Teo",
#      "name": "teo",
#      "password": "pass",
#      "api_token": "44444"
# }
# req = requests.request('PUT', add_patient_url, json=data1)
#
#
# data2 = {
#      "username": "Alejandro",
#      "name": "bobita",
#      "password": "pass",
#      "api_token": "44444"
# }
# req = requests.request('PUT', add_patient_url, json=data2)
#
#
# data3 = {
#      "username": "costelino",
#      "name": "costelino",
#      "password": "pass",
#      "api_token": "44444"
# }
# req = requests.request('PUT', add_patient_url, json=data3)
#
# data4 = {
#      "username": "Carton",
#      "name": "qeqeqe",
#      "password": "pass",
#      "api_token": "44444"
# }
# req = requests.request('PUT', add_patient_url, json=data4)
#
# data4 = {
#      "username": "Czd",
#      "name": "riciard",
#      "password": "pass",
#      "api_token": "44444"
# }
# req = requests.request('PUT', add_patient_url, json=data4)



def retrieve(request):
    header = {
        'access_token': request.user.userprofile.access_token}
    retrieve_url = "http://192.168.50.4:2002/api/general/physicians/"\
                   +str(request.user.userprofile.id_big_database)+"/retrieve"
    retrieve_req = requests.request('GET', url=retrieve_url, headers=header)
    print (retrieve_req.text)

    text = retrieve_req.text
    text = text.split(',')

    dict = {}
    if len(text) > 1:
        for x in range(len(text)):
            if x % 2 == 0:
                # print (text[x].split(':')[1])
                # print (text[x+1].split(':')[1].split('"')[1])
                dict[text[x].split(':')[1]] = text[x+1].split(':')[1].split('"')[1]
                # print ("+++++++++++++++++")
    return dict


@login_required
def patient_list(request):
    context = ({
                'patients': retrieve(request)
                })
    return render(request, "patient_list_landing.html", context=context)


@login_required
def patient_details(request, pk):
    print("++++++++++++++++++++++++")
    print(pk)
    print("++++++++++++++++++++++++")
    l = [1, 2]
    context = ({'pk': pk,
                'patients': retrieve(request)
               })
    return render(request, "patient_list_details.html", context=context)


@login_required
def add_patient(request):

    if request.method == 'POST':
        form = PatientForm(request.POST)
        if form.is_valid():
            name = form.cleaned_data['name']
            username = form.cleaned_data['username']
            password = form.cleaned_data['password']
            api_token = form.cleaned_data['api_token']
            print("name "+name)
            print("username "+username)
            print(password)
            print(api_token)
            print(request.user.username)
            add_patient_url = "http://192.168.50.4:2002/api/accounts/patients?token="+request.user.username.upper()
            print (add_patient_url)
            data = {
                "username": username,
                "name": name,
                "password": password,
                "api_token": api_token
            }
            print(data)
            req = requests.request('PUT', add_patient_url, json=data)
            print(req)
            if req.status_code == 500:
                return render(request, "wrong.html")
            else:
                return redirect('api:list')

    else:
        form = PatientForm()

    return render(request, "add_patient.html", {'form': form})


def treatment(request):
    if request.method == 'POST':
        time = request.POST.get('time', False)
        date = request.POST.get('date', False)
        option = request.POST.get('optradio', False)
        print (time)
        print (date)
        print (option)
        return redirect('api:list')
    return render(request, "treatment.html")

@user_passes_test(lambda u: not u.is_active)
def login(request):
    if request.method == 'POST':
        form = LoginForm(request.POST)
        if form.is_valid():
            username = form.cleaned_data['username']
            password = form.cleaned_data['password']
            print (username)
            print (password)
            data = {
                "username": username,
                "password": password
            }
            print (data)
            login_url = "http://192.168.50.4:2002/api/accounts/login"
            login_doctor_req = requests.request('POST', login_url, json=data)

            if login_doctor_req.status_code == 500:
                return render(request, "wrong.html")
            else:
                # a = UserProfile.objects.create_object(username=username, password=password)
                print (login_doctor_req.text)

                first = login_doctor_req.text.split(',')[0][10:-1]
                second = login_doctor_req.text.split(',')[1][5:-1]
                #
                user = User(
                    username=username, password=password)
                user.save()
                a = UserProfile(user=user, access_token=first, id_big_database=int(second))
                a.save()
                print (first)
                print (second)

                # user = authenticate(username=username, password=password)
                auth_login(request, user)

                return redirect('api:list')
    else:
        form = LoginForm()
    return render(request, "login.html", {'form': form})


@login_required
def user_logout(request):

    if request.method == 'POST':

        if request.POST.get("yes_button"):
            print ("am ajuns aiciaaaa")
            user = request.user
            logout(request)
            user.delete()
            return redirect('api:login')

        elif request.POST.get("no_button"):
            print ("nuuuuu")
            return redirect('api:list')

    return render(request, "logout.html")

@user_passes_test(lambda u: not u.is_active)
def register(request):
    if request.method == 'POST':
        form = RegisterForm(request.POST)
        if form.is_valid():
            username = form.cleaned_data['username']
            name = form.cleaned_data['name']
            password = form.cleaned_data['password']
            api_token = form.cleaned_data['api_token']
            email = form.cleaned_data['email']
            creation_token = form.cleaned_data['creation_token']
            print (username)
            print (name)
            print (password)
            print (api_token)
            print (email)
            print (creation_token)
            print("name " + name)
            print("username " + username)
            print(password)
            print(api_token)
            data = {
                "username": username,
                "name": name,
                "password": password,
                "api_token": api_token,
                "email": email,
                "creation_token": creation_token
            }
            print (data)
            create_physician_url = "http://192.168.50.4:2002/api/accounts/physicians"
            create_physician_req = requests.request('PUT', url=create_physician_url, json=data)
            print (data)
            print (create_physician_req)
            print (create_physician_req.text)

            if create_physician_req.status_code == 500:
                return render(request, "wrong.html")
            else:
                return redirect('api:login')
    else:
        form = RegisterForm()
    return render(request, "register.html", {'form': form})


