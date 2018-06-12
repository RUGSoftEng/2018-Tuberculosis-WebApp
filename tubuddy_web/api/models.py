from django.db import models
from django.contrib.auth.models import User


class UserProfile(models.Model):

    user = models.OneToOneField(User, on_delete=models.CASCADE)
    # name = models.CharField(max_length=100)
    # api_token = models.CharField(max_length=300)
    access_token = models.CharField(max_length=300)
    id_big_database = models.IntegerField()

    def __str__(self):
        return self.user.username
