from django import forms


class PatientForm(forms.Form):
    name = forms.CharField(label='Name', max_length=100)
    username = forms.CharField(label='Username', max_length=100)
    password = forms.CharField(label='Password', max_length=100)
    api_token = forms.CharField(label='Api_token', max_length=300)