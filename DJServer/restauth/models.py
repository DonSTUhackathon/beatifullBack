from os import name
from django.db import models
from django.contrib.auth.models import User
from django.db.models.base import CASCADE

# Create your models here.


class Profile(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    description = models.TextField()

class room(models.Model):
    user1 = models.ForeignKey(Profile, on_delete=CASCADE, related_name="user1")
    user2 = models.ForeignKey(Profile, on_delete=CASCADE, related_name="user2")

    
        



    





