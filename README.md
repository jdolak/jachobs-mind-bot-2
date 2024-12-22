# Jachob's Mind 2

## Into

Jachob's Mind is a discord bot that fulfills niches that I have desired in my discord servers. 
The original Jachob's Mind was a bot that I wrote in highschool to keep track of my schools changing bell schedule.
This new bot switches from Python to Go and acts as a bookkeeper for my group of friend's activites in and out of discord.

## Features

### Credit

Credit is a virtual "currency" that is a flexible, zero sum, and peer modifiable.  
It is purposely generic to fit any theme of the discord server. Credit can represent points in a game, money in a store, or reputation in a forum.

Credit can also be tied to certain bot and server features, allowing users to unlock or purchase different abilities.  

#### Commands

To view a user's current credit amount:
```
\credit <user>
```
To increase or decrease a user's credit by a certain amount:
```
\credit <user> <difference>
```
To view the top and bottom three credit holders:
```
\leaderboard
```


### Debt

While *credit* acts as a global amount that is each user can modify, *debt* represents a relationship between two users. While it can be flexible, the goal of *debt* to be an simple way of keeping track of money owed IRL.

Whether its spotting someone for lunch or splitting expenses on a trip, *debt* makes it easy with simple commands, bill splitting, and generating venmo links with the exact amount owed.

This makes no guarantees about getting your money back, just offers a convenient way to bookkeep.

#### Commands

To split a bill between muliple people, including yourself:
```
\split ...
```
To create a debt for youself (e.g "You owe user $10"):
```
\owes ...
```

To create a debt for someone to you (e.g "user owes you $10"):
```
\loan ...
```

To decrease you debt to someone (e.g "you paid $5 back to user"):
```
\paid ...
```

To decrease the debt someone has to you (e.g. "user paid you back $5"):
```
\recieved ...
```

Generate venmo links for everyone you owe money to, preset with the amount you owe them:
```
\link ...
```

To register your venmo account so other people can pay you through a link
```
\register ...
```

