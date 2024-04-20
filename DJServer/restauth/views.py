from django.contrib.admin.options import settings
from django.db.models.fields import return_None
from rest_framework.compat import requests
from allauth.socialaccount.providers.google.views import GoogleOAuth2Adapter
from allauth.socialaccount.providers.oauth2.client import OAuth2Client
from dj_rest_auth.registration.views import SocialLoginView
from DJServer.settings import GOOGLE_REDIRECT_URL
from django.http import HttpResponse
import requests
from django.contrib.auth.mixins import LoginRequiredMixin
from django.views.generic import RedirectView
from rest_framework import generics
from rest_framework.views import APIView

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
        r = requests.post("http://localhost:8080/auth/google/cb/", data=payload)
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



class MathMe(APIView):
    pass




        










