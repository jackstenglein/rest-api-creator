from django.shortcuts import render, redirect
from django.http import HttpResponse
from django.forms import formset_factory
from django.contrib.auth import login, authenticate
from django.contrib.auth.forms import UserCreationForm

from .models import Object
from .forms import ObjectForm, AttributeForm

# Create your views here.
def signup(request):
    if request.method == 'POST':
        form = UserCreationForm(request.POST)
        if form.is_valid():
            form.save(commit=False)
            username = form.cleaned_data.get('username')
            raw_password = form.cleaned_data.get('password1')
            user = authenticate(username=username, password=raw_password)
            login(request, user)
            return redirect('main_app:get_objects')
    else:
        form = UserCreationForm()
    return render(request, 'main_app/signup.html', {'form': form})


def get_objects(request):
    objects = Object.objects.all()
    print(objects)
    return render(request, 'main_app/objects.html', {'objects': objects})


def create_object(request):
    '''
    if request.method == "POST":
        print(request.POST)
        return redirect('get_objects', request)
    else:
        return render(request, 'main_app/create_object.html')
    '''
    ObjectFormSet = formset_factory(ObjectForm)
    AttributeFormSet = formset_factory(AttributeForm)
    if request.method == 'POST':
        object_formset = ObjectFormSet(request.POST, request.FILES, prefix='object')
        attribute_formset = AttributeFormSet(request.POST, request.FILES, prefix='attributes')
        if object_formset.is_valid() and attribute_formset.is_valid():
            print("Everything is valid")
            new_model = object_formset[0].save(commit=False)
            print(new_model)
            for form in attribute_formset:
                new_attribute = form.save(commit=False)
                print(new_attribute)
            pass
    else:
        object_formset = ObjectFormSet(prefix='object')
        attribute_formset = AttributeFormSet(prefix='attributes')
    return render(request, 'main_app/create_object2.html', {
        'object_formset': object_formset,
        'attribute_formset': attribute_formset,
    })
