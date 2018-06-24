from django import forms


class PatientForm(forms.Form):
    name = forms.CharField(widget=forms.TextInput(attrs={'class': 'form-control'}), label='Name', max_length=100)
    username = forms.CharField(widget=forms.TextInput(attrs={'class': 'form-control'}), label='Username', max_length=100)
    password = forms.CharField(widget=forms.PasswordInput(attrs={'class': 'form-control'}), label='Password', max_length=100)
    # api_token = forms.CharField(widget=forms.TextInput(attrs={'class': 'form-control'}), label='Api_token', max_length=300)


class LoginForm(forms.Form):
    username = forms.CharField(widget=forms.TextInput(attrs={'class': 'form-control'}), label='Username', max_length=100)
    password = forms.CharField(widget=forms.PasswordInput(attrs={'class': 'form-control'}), label='Password',
                               max_length=100)


class RegisterForm(forms.Form):
    username = forms.CharField(widget=forms.TextInput(attrs={'class': 'form-control'}), label='Username',
                               max_length=100)
    name = forms.CharField(widget=forms.TextInput(attrs={'class': 'form-control'}), label='Name', max_length=100)
    password = forms.CharField(widget=forms.PasswordInput(attrs={'class': 'form-control'}), label='Password',
                               max_length=100)
    # api_token = forms.CharField(widget=forms.TextInput(attrs={'class': 'form-control'}), label='Api_token', max_length=100)
    email = forms.EmailField(widget=forms.TextInput(attrs={'class': 'form-control'}), label='Email',
                                max_length=100)
    creation_token = forms.CharField(widget=forms.TextInput(attrs={'class': 'form-control'}), label='Creation_token',
                                max_length=100)