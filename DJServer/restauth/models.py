from os import name
from django.db import models
from django.contrib.auth.models import User
from django.db.models.base import CASCADE
from django.db.models.fields import CharField
from rest_framework import serializers

# Create your models here.
class Profile(models.Model):
    id = models.OneToOneField(User, on_delete=models.CASCADE, primary_key=True)
    description = models.TextField()
    image_path = models.CharField(max_length=100)



class Room(models.Model):
    user1 = models.ForeignKey(Profile, on_delete=CASCADE, related_name="user1")
    user2 = models.ForeignKey(Profile, on_delete=CASCADE, related_name="user2")
    

class Meeting(models.Model):
    desc = CharField(max_length=100)
    long = models.DecimalField(max_digits=9, decimal_places=6)
    lat = models.DecimalField(max_digits=9, decimal_places=6)
    time = models.TimeField()
    room = models.ForeignKey(Room, on_delete=CASCADE)

    
        
class Tag(models.Model):
    name = models.CharField(max_length=100)
    users = models.ManyToManyField(Profile)



# serializers
class RoomS(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Room
        fields = ['user1', 'user2']

class MeetingsS(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Meeting
        fields = ['desc', 'long', 'lat', 'time', 'room']

class ProfileS(serializers.HyperlinkedModelSerializer):
    id = serializers.HyperlinkedRelatedField(view_name="user-detail", queryset=User.objects.all())
    class Meta:
        model = Profile
        fields=['id', 'description', 'image_path']

class UserS(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = User
        fields = ['username', 'email']








