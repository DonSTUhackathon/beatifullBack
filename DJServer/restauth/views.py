import random
import requests
from itertools import chain

from DJServer.settings import GOOGLE_REDIRECT_URL

from django.http import Http404
from allauth.socialaccount.providers.google.views import GoogleOAuth2Adapter
from allauth.socialaccount.providers.oauth2.client import OAuth2Client
from dj_rest_auth.registration.views import Response, SocialLoginView

from django.http import HttpResponse
from django.contrib.auth.mixins import LoginRequiredMixin
from django.views.generic import RedirectView


from rest_framework import generics
from rest_framework.authentication import SessionAuthentication, BasicAuthentication
from rest_framework.permissions import IsAdminUser, IsAuthenticated
from rest_framework.views import APIView
from rest_framework.compat import requests  
from rest_framework import serializers, viewsets

from dj_rest_auth import jwt_auth
from restauth.models import *

from rest_framework_simplejwt.tokens import RefreshToken

def get_tokens_for_user(user):
    refresh = RefreshToken.for_user(user)

    return {
        'refresh': str(refresh),
        'access': str(refresh.access_token),
    }

def profile(user):
    return Profile.objects.filter(id=user)[0]

class GoogleLoginView(SocialLoginView):# Custom adapter is created because obsolete URLs are used inside django-allauth library
    class GoogleAdapter(GoogleOAuth2Adapter):
        access_token_url = "https://oauth2.googleapis.com/token"
        authorize_url = "https://accounts.google.com/o/oauth2/v2/auth"
        profile_url = "https://www.googleapis.com/oauth2/v2/userinfo"    
    adapter_class = GoogleAdapter
    callback_url = GOOGLE_REDIRECT_URL 
    client_class = OAuth2Client

    def get(self, request):
        try:
            code = request.GET['code']
            scope = request.GET['scope']
            print(code, scope)
        except KeyError:
            return HttpResponse("No credits in request!")
        
        payload = {
                "code":code,

                }
        r = requests.post(GOOGLE_REDIRECT_URL, data=payload, verify=False)
        print(request.user)

        return HttpResponse(f"Code:{code}\nScope:{scope}, {r.text} ")


class UserRedirectView(LoginRequiredMixin, RedirectView):
    """
    This view is needed by the dj-rest-auth-library in order to work the google login. It's a bug.
    """
    permanent = False

    def get_redirect_url(self):
        return "redirect-url"

def testView(request):
    return HttpResponse("Test resp") 

class ProfileViews(viewsets.ModelViewSet):
    authentication_classes=[SessionAuthentication, BasicAuthentication]
    #permission_classes = [IsAuthenticated]
    queryset = Profile.objects.all()
    serializer_class = ProfileS

class UserViews(viewsets.ModelViewSet):
    authentication_classes=[SessionAuthentication, BasicAuthentication]
    #permission_classes = [IsAuthenticated]
    queryset = User.objects.all()
    serializer_class = UserS

class TagViews(viewsets.ModelViewSet):
    authentication_classes=[SessionAuthentication, BasicAuthentication]
    #permission_classes = [IsAuthenticated]
    queryset = Tag.objects.all()
    serializer_class = TagS


class RoomViews(viewsets.ModelViewSet):
    authentication_classes=[SessionAuthentication, BasicAuthentication]
    #permission_classes = [IsAdminUser]
    queryset=Room.objects.all()
    serializer_class = RoomS

class MatchMe(APIView):

    authentication_classes=[SessionAuthentication, BasicAuthentication]
    #permission_classes = [IsAuthenticated]
    
    def get(self, request, format=None):
        print(request.user)
        q1 =  Room.objects.filter(user1=profile(request.user))
        q2 = Room.objects.filter(user2=profile(request.user))
        if q1 or q2:
            print("Has rooms//")
            abandoned = list(chain(q1, q2))
            choose = Profile.objects.exclude(abandoned)
        else:
            choose = Profile.objects.all()
        if not choose:
            return Http404 

        for p in choose:
            print(" ", p.id)
        
        chosed = random.choice(choose)
        print("Chosed", chosed.id)

        # Вернуть Room
        context = {'request':request}
        new_room = Room(profile(request.user), chosed)
        new_room.save()
        s = ProfileS(choosed, context=context)
        print(f"Creating room wtih {request.user} {chosed}!")
        return Response(s.data)

            

        




        










