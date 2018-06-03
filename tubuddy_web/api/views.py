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





def retrieve():
    # header = {
    #     'access_token': 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZCI6InBhc3MiLCJ1c2VybmFtZSI6IkhvdXNlIn0.JcFAPE8Y9xlgP9hbef0qTO5tZnQLSZDKOqL3MgIliuM'}
    # retrieve_url = "http://192.168.50.4:2002/api/general/physicians/3/retrieve"
    # retrieve_req = requests.request('GET', url=retrieve_url, headers=header)
    # print (retrieve_req.text)
    #
    #
    # text = retrieve_req.text
    # # text = '[{"id":4,"name":"Andrei"},{"id":5,"name":"Cosmin"},{"id":6,"name":"Teodor"},{"id":7,"name":"Costelusu"},{"id":8,"name":"lelelelele"},{"id":9,"name":"bossu1"},{"id":10,"name":"bossu123"},{"id":11,"name":"bossu1234"},{"id":12,"name":"coiosu"}]'
    # text = text.split(',')
    #
    # dict = {}
    # if len(text) > 1:
    #     for x in range(len(text)):
    #         if x % 2 == 0:
    #             # print (text[x].split(':')[1])
    #             # print (text[x+1].split(':')[1].split('"')[1])
    #             dict[text[x].split(':')[1]] = text[x+1].split(':')[1].split('"')[1]
    #             # print ("+++++++++++++++++")
    # return dict
    return {}


@login_required
def patient_list(request):
    l = {'puli': "dadada",
         'castravetoi': "sa fie paceeee"
         }
    # print (login_doctor_req.text)
    # print (type(login_doctor_req.text))
    context = ({
                'patients': retrieve()
                })
    return render(request, "patient_list_landing.html", context=context)


@login_required
def patient_details(request, pk):
    print("++++++++++++++++++++++++")
    print(pk)
    print("++++++++++++++++++++++++")
    l = [1, 2]
    context = ({'pk': pk,
                'patients': retrieve()
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
            add_patient_url = "http://192.168.50.4:2002/api/accounts/patients?token=HOUSE"
            # add_patient_url = "http://192.168.50.4:2002/api/accounts/patients?token="
            # User.objects.get().username
            print (add_patient_url)
            import ipdb; ipdb.set_trace()
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
                return render(request, "muie.html")
            else:
                # return HttpResponseRedirect("/")
                context = ({
                        'patients': retrieve()
                        })
                return render(request, "patient_list_landing.html", context=context)

    else:
        form = PatientForm()

    return render(request, "add_patient.html", {'form': form})


def treatment(request):
    if request.method == 'POST':
        time = request.POST.get('time', False)
        date = request.POST.get('date', False)
        print (time)
        print (date)
        return render(request, "patient_list_landing.html")
        # return HttpResponseRedirect("/")
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
                return render(request, "muie.html")
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

                # import ipdb; ipdb.set_trace()
                return HttpResponseRedirect("/")
    else:
        form = LoginForm()
    return render(request, "login.html", {'form': form})


@login_required
def user_logout(request):
    # import ipdb; ipdb.set_trace()
    # user = request.user
    #
    # logout(request.user.delete())

    # import ipdb; ipdb.set_trace()

    if request.method == 'POST':

        if request.POST.get("yes_button"):
            print ("am ajuns aiciaaaa")
            user = request.user
            logout(request)
            user.delete()
            return HttpResponseRedirect("/")

        elif request.POST.get("no_button"):
            print ("nuuuuu")
            return HttpResponseRedirect("/")

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
            create_physician_url = "http://192.168.50.4:2002/api/accounts/physicians"
            create_physician_req = requests.request('PUT', url=create_physician_url, json=data)
            print (data)
            print (create_physician_req)
            print (create_physician_req.text)


            if create_physician_req.status_code == 500:
                return render(request, "muie.html")
            else:
                return HttpResponseRedirect("/")
    else:
        form = RegisterForm()
    return render(request, "register.html", {'form': form})



def test(request):
    # print(request.readlines())
    # import ipdb; ipdb.set_trace()
    if request.method == 'POST':
        import ipdb; ipdb.set_trace()
        name = request.POST.get('name', False)
        email = request.POST.get('email', False)
        phone = request.POST.get('phone', False)
        print (name)
        print (email)
        print (phone)
        # return HttpResponseRedirect("http://www.google.com")
        return render(request, "patient_list_landing.html")
    return render(request, "test.html")

def test2(request):
    # print(request.readlines())
    # import ipdb; ipdb.set_trace()
    if request.method == 'POST':
        name = request.POST.get('name', False)
        email = request.POST.get('email', False)
        phone = request.POST.get('phone', False)
        print (name)
        print (email)
        print (phone)
        # return HttpResponseRedirect("http://www.google.com")
        return render(request, "patient_list_landing.html")
    return render(request, "test2.html")