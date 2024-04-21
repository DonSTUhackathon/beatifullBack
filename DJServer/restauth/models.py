from os import name
from django.db import models
from django.contrib.auth.models import User
from django.db.models import query
from django.db.models.base import CASCADE
from django.db.models.expressions import fields
from django.db.models.fields import CharField
from rest_framework import serializers

class Tag(models.Model):
    name = models.CharField(max_length=100)


# Create your models here.
class Profile(models.Model):
    id = models.OneToOneField(User, on_delete=models.CASCADE, primary_key=True)
    description = models.TextField()
    image_path = models.CharField(max_length=100)
    tags = models.ManyToManyField(Tag)



class Room(models.Model):
    user1 = models.ForeignKey(Profile, on_delete=CASCADE, related_name="user1") # Запрашиваемый
    user2 = models.ForeignKey(Profile, on_delete=CASCADE, related_name="user2")

    

class Meeting(models.Model):
    desc = CharField(max_length=100)
    time = models.TimeField()
    room = models.ForeignKey(Room, on_delete=CASCADE)
    accepted = models.BooleanField()


    
        



class TagS(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Tag
        fields = ['name']

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
    tags = TagS(many=True) 
    class Meta:
        model = Profile
        fields=['id', 'description', 'image_path', 'tags']

class UserS(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = User
        fields = ['username', 'email']








