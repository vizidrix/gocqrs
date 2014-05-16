Because having a usecase called "customer" is really really weird. That 
sounds like a domain name, and not a usecase name, is what I mean.

I can compare with the problem I ran into in my own app. I used to have 
this:
/user/rickard/

from which I could do two commands: changepassword and resetpassword. 
But that becomes REALLY funky when you try to do authorization and 
linking properly, since they are part of two different client usecases. 
So instead I changed this to:
/account/ (my own user is implied)
/administration/server/user/rickard/

i.e. the usecases are "account management" and "user administration", 
with changepassword and resetpassword respectively. With that change it 
became much easier to handle, on the server and client. So, again, 
"customer" to me doesn't sound like a usecase. It sounds like exposing 
the domain model, and that breaks down fantastically when you try to do 
HATEOAS properly.

/Rickard



http://www.claassen.net/geek/blog/2011/04/http-cqrs-restrpc.html


After thinking about this for a bit longer, I don't like my proposed POST anymore. The issue is that it strongly binds the implementation to a body type that can represent the command.  That makes it impossible to use simpler bodies like CSVs without wrapping a document arount them. So, I'm leaning towards POST:users/{id}/{command}. Alternatively, there could be a ?command= query parameter.  Either way, the command distinction is removed from the request body.
Also, there should be a way to query for all available views and commands.  Something like, GET:users/views and GET:users/commands. That's assuming, {id} can't be 'views' or 'commands' of course.



GET == View

api/{version}/views/
- List of all general views

api/{version}/views/{viewname}/
- Returns contents of gneral view

api/{version}/{aggregate}/
- Default aggregate view?

api/{version}/{aggregate}/commands/
- List of available commands

api/{version}/{aggregate}/commands/{commandname}/
- Details of command parameters

api/{version}/{aggregate}/events/
- List of available events

api/{version}/{aggregate}/events/{eventname}/
- Details of event parameters

api/{version}/{aggregate}/views/
- List of views for this aggregate?

api/{version}/{aggregate}/views/{viewname}/
- Returns contents of aggregate view

api/{version}/{aggregate}/{id}/
- Default aggregate instance view?

api/{version}/{aggregate}/{id}/commands/
- List of available commands

api/{version}/{aggregate}/{id}/commands/{commandname}/
- Details of command parameters

api/{version}/{aggregate}/{id}/events/
- List of available events

api/{version}/{aggregate}/{id}/events/{eventname}/
- Details of event parameters

api/{version}/{aggregate}/{id}/views/
- List of views for this aggregate instance?

api/{version}/{aggregate}/{id}/views/{viewname}/
- Returns contents of aggregate view


POST == Command

api/{version}/{aggregate}/{commandname}/

api/{version}/{aggregate}/{id}/{commandname}/

PUT == Event

api/{version}/{aggregate}/{eventname}/

api/{version}/{aggregate}/{id}/{eventname}/


