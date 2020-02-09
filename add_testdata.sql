PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;
DELETE FROM notereply WHERE 1=1;
DELETE FROM note WHERE 1=1;
DELETE FROM file WHERE 1=1;
DELETE FROM user WHERE user_id > 2;

INSERT INTO user (username, password) VALUES ('robdelacruz', '$2a$10$QBKdo66QfkyqNczexwGFwul3731pQ970B96Bn1hgmvXLBu.LaJhFK'); -- password is '123'
INSERT INTO user (username, password) VALUES ('lky', '');

INSERT INTO note (title, body, createdt, user_id) VALUES ('Aimee Teagarden', 'All about Aimee Teagarden Hallmark show', '2019-12-01T14:00:00+08:00', 3);
INSERT INTO note (title, body, createdt, user_id) VALUES ('Emma Fielding', 'All about Emma Fielding Hallmark show', '2019-12-02T14:00:00+08:00', 3);
INSERT INTO note (title, body, createdt, user_id) VALUES ('Mystery 101', 'All about Mystery 101 Hallmark show', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO note (title, body, createdt, user_id) VALUES ('test note', 'test note 1', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO note (title, body, createdt, user_id) VALUES ('test note 2', 'test note 2', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO note (title, body, createdt, user_id) VALUES ('test note 3', 'test note 3', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO note (title, body, createdt, user_id) VALUES ('markdown test', '# Gettysburg Address

*Versions*

- Bliss Copy
- Nicolay Copy
- Hay Copy
- Everett Copy
- Bancroft Copy

### Related Links

[Robert Todd Lincoln''s "Gettysburg Story"](https://quod.lib.umich.edu/j/jala/2629860.0038.103/--robert-todd-lincolns-gettysburg-story?rgn=main;view=fulltext) (JALA)
[Who stole the Gettysburg Address?](https://quod.lib.umich.edu/j/jala/2629860.0024.203/--who-stole-the-gettysburg-address?rgn=main;view=fulltext) (JALA)

---

Four score and seven years ago our fathers brought forth on this continent, a new nation, conceived in Liberty, and dedicated to the proposition that all men are created equal.

Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting place for those who here gave their lives that that nation might live. It is altogether fitting and proper that we should do this.

But, in a larger sense, we can not dedicate -- we can not consecrate -- we can not hallow -- this ground. The brave men, living and dead, who struggled here, have consecrated it, far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us -- that from these honored dead we take increased devotion to that cause for which they gave the last full measure of devotion -- that we here highly resolve that these dead shall not have died in vain -- that this nation, under God, shall have a new birth of freedom -- and that government of the people, by the people, for the people, shall not perish from the earth.

Abraham Lincoln
November 19, 1863

![Soldiers National Cemetery](http://www.abrahamlincolnonline.org/lincoln/sites/gettycem.jpg)

[source](http://www.abrahamlincolnonline.org/lincoln/speeches/gettysburg.htm)
', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO note (title, body, createdt, user_id) VALUES ('code test', '## Hello, World

Code for Hello, World:

    #include <stdio.h>

    int main() {
        printf("Hello, World!\n");
    }
', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO note (title, body, createdt, user_id) VALUES ('Lee Kuan Yew Quotes', '### Lee Kuan Yew Quotes:

>"If there was one formula for our success,it was that we were constantly studying how to make things work,or how to make them work better."

>"I’m very determined. If I decide that something is worth doing, then I’ll put my heart and soul to it. The whole ground can be against me, but if I know it is right, I’ll do it."

>"I always tried to be correct, not politically correct."

', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (5, 'first comment!', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (5, 'second comment!', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (5, 'third comment!', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (6, 'a comment!', '2019-12-05T14:00:00+08:00', 3);
INSERT INTO notereply (note_id, replybody, createdt, user_id) VALUES (6, 'another comment!', '2019-12-05T14:00:00+08:00', 3);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Comanche',  
'You need to straighten your posture and suck in your gut
You need to pull back your shoulders and tighten your butt
Come comanche, comanche, comanche, come oh
If you want to have cities, you''ve got to build roads

You need to find some new feathers and buy some new clothes
Just get rid of the antlers and lighten your load
Come comanche, comanche, comanche, come oh
If you want to have cities, you''ve got to build roads

You need to straighten your posture and suck in your gut
You need to pull back your shoulders and tighten your butt
Come comanche, comanche, comanche, come oh
If you want to have cities
If you want to have cities
If you want to have cities, you''ve got to build roads' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Ruby Sees All',  
'Do you think she''s swimming in your lies?
Do you think it''s all just murky green?
Don''t you think that she would realize?
Yeah, do you think that she has never seen?

[Chorus]
''Cause when the seaweed sinks
And the sun gets low
When the waves retire
To the darkness below
I know
I know Ruby sees all
Whoa, I know
I know Ruby sees

I can feel the pressure building high
You should see you''re headed for a storm
Don''t you see it building in the sky?
Don''t you think it''s time to swim to shore?

[Chorus]'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Up So Close',  
'Up so close
I never get to see your face
Microscope
I might as well be out in space

Up so close
I never get to see the view
Down your throat
I''m never sure that it''s still you

Up your nose
Down to your toes
In your mouth
Way down south

Up so close
It seems I only think of you
Up so close
I never see the sky so blue

I only wanted to be sure
That what it was was really pure
I put my face down in the cake
My feet were flailing in a lake

Up so close
I never get to see your face
Microscope
I might as well be out in space

Up so close
I never get to see you
Microscope
I''m never sure if it''s still you'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Pentagram',  
'Your pentagram is down below our floor
Your naked body shimmers in the night
Dancing and chanting in a sacrificial rite
Your feet are dry with the ashes from dead babies
Who have passed the test
Just like all the rest
But never really understood
The reasons why they took it
In the first place
Ahh, in the first place

Your feasty eyes won''t make me fall apart
Your turquoise and silver won''t weaken this old heart
Yeah, dancing and chanting in a sacrificial rite
I fell to the ground on a windy, windy night

Well I have passed the test
Just like all the rest
But never really understood
The reasons why I took it
In the first place
Ahh, in the first place'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Jolene',  
'Well, Jolene unlocked the thick breezeway door
Like she''d done one hundred times before
Jolene smoothed her dark hair in the mirror
She folded the towel carefully and put it back in place

Yeah, I want to pull you down into bed
I want to cast your face in lead

[Chorus]
But every time I pull you close
Push my face into your hair
Cream rinse and tobacco smoke
That sickly scent is always, always there
Yeah, yeah

Jolene heard her father''s uneven snores
Right then, she knew there must be something more
Jolene heard the singing in the forest
She opened the door quietly
And stepped into the night

Yeah, I want to throw you out into space
I want to do whatever it takes

[Chorus]

Get down!
Get down, down!
Get up!
Get down!
Please, get down!
Get down!
Get down!
Get down!
Yeah, all right!
Yeah
That''s great, that''s great
Oh yeah
Get down!
Yeah
Oh yeah'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Haze of Love',  
'It''s 3 o''clock in the morning
Or maybe it''s 4
I''m thinking of you
Wondering what I should do
But I''m finally cutting
Through this haze

It''s 4 o''clock in the morning
Or maybe it''s 5
I think I''m alive
And I think I''ll survive
And I''m finally cutting
Through this haze of love
Haze of love
For days and days
I''m in a haze of love

Yea, you don''t love me
Like I love you
Although you pretend
I can see this will end
I''m finally cutting through
This haze of love
Haze of love
For days and days and days
I''m in a haze of love

It''s 5 o''clock in the morning
Or maybe it''s 6
I am sick of your lies
I am sick of your tricks
I am finally cutting through
This haze of love
Haze of love
For days and days and days
For days and days and days
For days and days and days
I''m in a haze of love'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'You Part the Waters',  
'You part the waters
The same ones that I''m drowning in
You lead your casual slaughters
And I''m the one who helps you win

You''ve got your grand piano
You don''t even play piano
I''m the one who plays piano
You don''t even play piano

You part the waters
The same ones that I''m thirsty for
You invite your friends to tea
But when it''s me you lock the door
You''ve got your credit cards
And you thank your lucky stars
But don''t forget the ones who foot the bill

You''ve got your grand piano
And you don''t even play piano
I''m the one who plays piano
You don''t even play piano
But you part the waters'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Is This Love?',  
'I can''t believe it''s true
I can''t believe that you
Don''t want me anymore
You''re with him
And you don''t even know
That I''ve been dying all day long
And singing sad, sad songs
And wishing you were gone

Is this love?
Is this love?
Is this love?
Or should I close the door?

My eyes are burning in my head
And seeing only red
And wishing you were dead

Is this love?
Is this love?
Is this love?
Is this love?
Is this love?
Or should I
Or should I, should I
Should I close the door?
Should I close the door? Ah, I fooled myself

Is this love?
Is this love?
Is this love?
Is this love?
Is this love?
Or should I
Or should I
Or should I, should I
Should I close the door?'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Jesus Wrote a Blank Check',  
'Jesus wrote a blank check
One I haven''t cashed quite yet
I hope I got a little more time
I hope it''s not the end of the line
Yeah, Jesus wrote a blank check
Ah, one I haven''t cashed yet, all right

Well, if I had to choose a number
I''d want it to be number one
I don''t want to be number two
Yeah, I don''t want to be number four
Yeah, but I can hear a knock on the door
Jesus wrote a blank check, all right

If Jesus saw me dying
Would angels come a-flying down?
I hope I got a little more time
I hope somebody lends me a dime
Now, Jesus wrote a blank check
Ah, one I haven''t cashed yet, uh-huh

Still I build my towers high
I watch them pierce the blue, blue sky
Still I wallow in the mire
Still I burn this earthen fire

Still I build my towers high
I watch them pierce the blue, blue sky
Still I wallow in the mire
Still I burn this earthen fire
Still I burn this earthen fire
Still I burn this earthen fire
Still I burn this earthen fire

Ah, still I burn this earthen fire
Still I burn this earthen fire
Still I burn this earthen fire
Still I burn this earthen fire'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Rock ''N'' Roll Lifestyle',  
'Well, your CD collection looks shiny and costly
How much did you pay for your Bad Moto Guzzi?
And how much did you spend on your black leather jacket
Is it you or your parents in this income tax bracket?
Now tickets to concerts
And drinking at clubs
Sometimes for music that you haven''t even heard of
And how much did you pay for your rock and roll t-shirt
That proves you were there
That you heard of them first?

[Chorus]
Now, how do you afford your rock ''n'' roll lifestyle?
How do you afford your rock ''n'' roll lifestyle?
How do you afford your rock ''n'' roll lifestyle?
Tell me

How much did you pay for the chunk of his guitar?
The one he ruthlessly smashed at the end of the show
And how much will he pay for a brand new guitar?
One which he''ll ruthlessly smash at the end of another show
And how long will the workers keep building him new ones?
As long as their soda cans are red, white, and blue ones
And how long will the workers keep building him new ones?
As long as their soda cans are red, white, and blue ones

Aging black leather
And hospital bills
And tattoo removal
And dozens of pills
Your liver pays dearly now for youthful magic moments
But rock on completely with some brand new components

[Chorus]

Excess ain''t rebellion
You drinkin'' what they''re sellin''
Your self-destruction doesn''t hurt them
Your chaos won''t convert them
They''re so happy to rebuild it
You''ll never really kill it
Excess ain''t rebellion
You drinkin'' what they''re sellin?
Excess ain''t rebellion
You drinkin'', you''re drinking what they''re... sellin'''
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'I Bombed Korea',  
'I bombed Korea every night
My engines sang into the salty sky
I didn''t know if I would live or die
I bombed Korea every night

I bombed Korea every night
I bombed Korea every night
Red flowers bursting down below us
Those people didn''t even know us
We didn''t know if we would live or die
We didn''t know if it was wrong or right
I bombed Korea every night

And so I sit here at this bar
I''m not a hero, I''m not a movie star
I''ve got my beer, I''ve got my stories to tell
But they won''t tell you what it''s like in hell
Red flowers bursting down below us
Those people didn''t even know us
We didn''t know if we would live or die
We didn''t know if it was wrong or right
We didn''t know if we would live or die
I bombed Korea every night'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Mr. Mastodon Farm',  
'Birds fall from the window ledge above mine
Then they flap their wings at the last second

You see, birds fall from the window ledge above mine
Then they flap their wings at the last second

I can see their dead weight
Just dropping like stones
Or small loaves of bread
Past my window all the time

But unless I get up
Walk across the room
And peer down below
I don''t see their last-second curves
Toward a horizontal flight
All these birds just falling from the ledge like stones

Now due to a construct in my mind
That makes their falling and their flight
Symbolic of my entire existence
It becomes important for me
To get up and see
Their last-second curves toward flight

It''s almost as if my life would fall
Unless I see their ascent

Mr. Mastodon Farm
Mr. Mastodon Farm
Cuts swatches out of all material

Mr. Mastodon Farm
Mr. Mastodon Farm
Cuts swatches out of all material'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Ain''t No Good',  
'She''s gonna hand you a red-headed Gabriel
Coming from the bar in a plastic tie
He''s gonna swing from the tree of life
He''s gonna try to sell you on a great big lie

But when you speak to her, her eyes light up
The music spills right into your cup
The minstrels play and the waitress brings ice
There are pies on a carousel, have a slice
But watch out, she ain''t no good for you

He''s gonna spin like the tractor pull
She''ll sit back when he tells his tale
He''s gonna yell when he drinks his beer
She''ll sit back and drink ginger ale

But when you speak to her, her eyes light up
The music spills right into your cup
It''s so abrupt and it''s so concise
There are pies on a carousel, have a slice
But watch out, she ain''t no good for you, I say
Watch out, she ain''t no good for you

She''d like to put you in her zoo
Right between the canaries and the cockatoos
She''ll pull out your feathers for her brand new hat
And when she''s done that, she''ll feed you to her cat
So watch out, she ain''t no good for you
Watch out, she ain''t no good for you
Watch out, she ain''t no good for you'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Frank Sinatra',  
'[Hook]
We know of an ancient radiation
That haunts dismembered constellations
A faintly glimmering radio station
While Frank Sinatra sings Stormy Weather
The flies and spiders get along together
Cobwebs fall on an old skipping record

[Verse 1]
Beyond the suns that guard this roost
Beyond your flowers of flaming truths
Beyond your latest ad campaigns
An old man sits collecting stamps
In a room all filled with Chinese lamps
He saves what others throw away
He says that he''ll be rich some day

[Bridge]
We know of an ancient radiation
That haunts dismembered constellations
A faintly glimmering radio station'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'The Distance',  
'[Verse 1]
Reluctantly crouched at the starting line
Engines pumping and thumping in time
The green light flashes, the flags go up
Churning and burning, they yearn for the cup
They deftly maneuver and muscle for rank
Fuel burning fast on an empty tank
Reckless and wild, they pour through the turns
Their prowess is potent and secretly stern
As they speed through the finish, the flags go down
The fans get up and they get out of town
The arena is empty except for one man
Still driving and striving as fast as he can

[Pre-Chorus]
The sun has gone down and the moon has come up
And long ago somebody left with the cup
But he''s driving and striving and hugging the turns
And thinking of someone for whom he still burns

[Chorus]
He''s going the distance
He''s going for speed
She''s all alone
In her time of need
Because he''s racing and pacing and plotting the course
He''s fighting and biting and riding on his horse
He''s going the distance

[Verse 2]
No trophy, no flowers, no flashbulbs, no wine
He''s haunted by something he cannot define
Bowel-shaking earthquakes of doubt and remorse
Assail him, impale him with monster-truck force
In his mind, he''s still driving, still making the grade
She''s hoping in time that her memories will fade
Cause he''s racing and pacing and plotting the course
He''s fighting and biting and riding on his horse

[Pre-Chorus]
The sun has gone down and the moon has come up
And long ago somebody left with the cup
But he''s striving and driving and hugging the turns
And thinking of someone for whom he still burns

[Chorus]
Cause he''s going the distance
He''s going for speed
She''s all alone
In her time of need
Because he''s racing and pacing and plotting the course
He''s fighting and biting and riding on his horse
He''s racing and pacing and plotting the course
He''s fighting and biting and riding on his horse
He''s going the distance
He''s going for speed
He''s going the distance' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Friend is a Four Letter Word',  
'To me, coming from you
Friend is a four letter word
End is the only part of the word
That I heard
Call me morbid or absurd
But to me, coming from you
Friend is a four letter word

But to me, coming from you
Friend is a four letter word
End is the only part of the word
That I heard
Call me morbid or absurd
But to me, coming from you
Friend is a four letter word

When I go fishing for the words
I am wishing you would say to me
I''m really only praying
That the words you''ll soon be saying
Might betray the way you feel about me

But to me, coming from you
Friend is a four letter word'
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Open Book',  
'She''s writing, she''s writing
She''s writing a novel
She''s writing, she''s weaving
Conceiving a plot
It quickens, it thickens
You can''t put it down now
It takes you, it shakes you
It makes you lose your thought
But you''re caught in your own glory
You are believing your own stories
Writing your own headlines
Ignoring your own deadlines
But now you''ve gotta write them all again

You think she''s an open book
But you don''t know which page to turn to, do you?
You think she''s an open book
But you don''t know which page to turn to, do you?
Do you? Do you?

You want her, confront her
Just open your window
Unbolt it, unlock it
Unfasten your latch
You want it, confront it
Just open your window
All you really have to do is ask

But you''re caught in your own glory
You are believing your own stories
Timing your contractions
Inventing small contraptions
That roll across your polished hardwood floors

You think she''s an open book
But you don''t know which page to turn to, do you?
You think she''s an open book
But you don''t know which page to turn to, do you?
Do you? Do you?

You think she''s an open book
But you don''t know which page to turn to, do you?
Do you? Do you? Do you?' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Daria',  
'When you tried to kiss me
I only bit your tongue
When you tried to get me together
I only came undone
When you tried to tell me
The one for me was you
I was in your mattress back in 1982

Daria, I won''t be soothed
Daria yeah, I won''t be soothed
Over like smoothed
Over like milk, silk
A bedspread or a quilt
Icing on a cake
Or a serene translucent lake

Daria, Daria, yeah Daria
I won''t be soothed
I won''t be soothed

When you tried to tell me
Of all the love you had
I was cleaning oil from beaches
Seeing only what was bad
When you tried to feed me
I only shut my mouth
Food got on your apron
And you told me to get out

Daria, I won''t be soothed
Daria yeah, I won''t soothed
Over like smoothed
Over like milk, silk
A bedspread or a quilt
Icing on a cake
Or a serene translucent lake

Daria, Daria, Daria
Daria yeah, Daria, yeah Daria yeah
I won''t be soothed
I won''t be soothed' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Race Car Ya-Yas',  
'The land of race car ya-yas
The land where you can''t change lanes
The land where large, fuzzy dice
Still hang proudly
Like testicles from rear-view mirrors
The land of race car ya-yas
The land where you can''t change lanes
The land where large, fuzzy dice
Still hang proudly
Like testicles from rear-view mirrors

The land of race car ya-yas
Ya-yas

The land of race car ya-yas
The land of race car ya-yas
Race car ya-yas' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'I Will Survive',  
'[Verse 1]
At first I was afraid
I was petrified
I kept thinking I could never live
Without you by my side
But then I spent so many nights
Just thinking how you''ve done me wrong
I grew strong
I learned how to get along
And so you''re back
From outer space
I just walked in to find you
Here without that look upon your face
I should have changed my fucking lock
I would have made you leave your key
If I''d have known for just one second
You''d be back to bother me

[Pre-Chorus]
Oh now go
Walk out the door
Just turn around now
You''re not welcome anymore
Weren''t you the one
Who tried to break me with desire?
Did you think I''d crumble?
Did you think I''d lay down and die?

[Chorus]
Oh not I
I will survive
Yeah
As long as I know how to love
I know I''ll be alive
I''ve got all my life to live
I''ve got all my love to give
I will survive
I will survive
Yeah, yeah

[Verse 2]
It took all the strength I had
Just not to fall apart
I''m trying hard to mend the pieces
Of my broken heart
And I spent oh so many nights
Just feeling sorry for myself
I used to cry
But now I hold my head up high
And you see me
With somebody new
I''m not that stupid little person
Still in love with you
And so you thought you''d just drop by
And you expect me to be free
But now I''m saving all my lovin''
For someone who''s lovin'' me

[Pre-Chorus]
Oh now go
Walk out the door
Just turn around now
You''re not welcome anymore
Weren''t you the one
Who tried to break me with desire?
Did you think I''d crumble?
Did you think I''d lay down and die?

[Chorus]
Oh not I
I will survive
Yeah
As long as I know how to love
I know I''ll be alive
I''ve got all my life to live
I''ve got all my love to give
I will survive
I will survive
Yeah, yeah
Da da, da da, da da dada...
Da, da, da, da, dada dada dada dada
...
Oh no' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Stickshifts and Safetybelts',  
'[Chorus]
Stick shifts and safety belts
Bucket seats have all got to go
When we''re driving in the car
It makes my baby seem so far
I need you here with me
Not way over in a bucket seat
I need you to be here with me
Not way over in a bucket seat

But when we''re driving in my Malibu
It''s easy to get right next to you
I say, "Baby, scoot over, please"
And then she''s right there next to me
I need you here with me
And not way over in a bucket seat
I need you to be here with me
Not way over in a bucket seat

Well, a lot of good cars are Japanese
Yeah, but when we''re driving far
I need my baby
I need my baby
Next to me

[Chorus]' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Perhaps, Perhaps, Perhaps',  
'You won''t admit you love me
And so
How am I ever
To know
You only tell me
Perhaps, perhaps, perhaps

A million times I ask you
And then
I ask you over
Again
You only answer
Perhaps, perhaps, perhaps

If you can''t make your mind up
We''ll never get started
And I don''t wanna'' wind up
Being parted, broken hearted

So if you really love me
Say yes
But if you don''t, dear
Confess
And please don''t tell me
Perhaps, perhaps, perhaps

If you can''t make your mind up
We''ll never get started
And I don''t wanna'' wind up
Being parted, broken hearted

So if you really love me
Say yes
But if you don''t, dear
Confess
And please don''t tell me
Perhaps, perhaps, perhaps
Perhaps, perhaps, perhaps
Perhaps, perhaps, perhaps' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'It''s Coming Down',  
'[Chorus]
It''s coming down
It''s coming down

[Verse 1]
It''s raining outside
You''ve nowhere to hide
She''s asking you
Why you think it''s funny

[Chorus]
It''s coming down
It''s coming down

[Verse 2]
She''s leaving your house
She had to get out
She''s mad
And she''ll take her mattress with her

[Chorus]
It''s coming down
It''s coming down
It''s coming down
It''s coming down
It''s coming down
It''s coming down
It''s coming down
It''s coming down
It''s coming down

[Verse 3]
You lie on the floor
She''s slamming your door
She''s gone
And she''s wearing your red sweater

[Chorus]
It''s coming down
It''s coming down
It''s coming down
It''s coming down
It''s coming down
It''s coming down
It''s coming down
It''s coming down
It''s coming down' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Nugget',  
'Okay alright, uh no
This one, this one, this one

Heads of state, who ride and wrangle
Who look at your face, from more than one angle
Can cut you from their bloated budgets
Like sharpened knives through Chicken McNuggets

Now heads of state, who ride and wrangle
Who look at your face, from more than one angle
Can cut you from their bloated budgets
Like sharpened knives through chicken McNuggets

Shut the fuck up, no
Shut the fuck up
(Shut the fuck)
Right, right
Learn to buck up
(Shut the fuck)

Right, shut the fuck up
Hey ho
(Shut the fuck)
Now, now
Learn to buck up

(Oh)
One, two, one two three four
Alright

Now, nimble fingers that dance on numbers
Will eat your children and steal your thunder
While heavy torsos that heave and hurl
Will crunch like nuts in the mouth of squirrels

Now, nimble fingers that dance on numbers
Will eat your children and steal your thunder
While heavy torsos that heave and hurl
Will crunch like nuts in the mouth of squirrels

Shut the fuck up, no
Shut the fuck up
(Shut the fuck)
Right, now
Learn to buck up

(Shut the fuck)
Right, shut the fuck up
Hey ho, ya
(Shut the fuck)
Ya ya
Learn to buck up

Now, simple feet that flicker like fire
And burn like candles in smoky spires
Do more to turn, my joy to sadness
Than somber thoughts of burning planets

Now, clever feet that flicker like fire
And burn like candles in smoky spires
Do more to turn, my joy to sadness
Than somber thoughts of burning planets

(Shut the fuck)
Alright, okay I don''t
(Shut the fuck)
Wanna, I don''t wanna hear it
That''s right
(Shut the fuck)
Oh, okay I don''t wanna
(Shut the fuck)

I don''t wanna
(Shut the fuck)
Hey, ho, ya
(Shut the fuck)
I don''t wanna
I don''t wanna
(Yea, one two one two one)' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'She''ll Come Back to Me',  
'[Verse 1]
Last night I said to her
I didn''t want to live inside a lie
If she wants him
More than she wants me
Let this be

[Chorus]
She''ll come back to me
She''ll come back to me
She''ll come back

[Verse 2]
All day I wait and wait
To hear her footsteps on my walkway
She never came
She never even called

[Chorus]
She''ll come back to me
She''ll come back to me
She''ll come back

[Bridge]
Somehow I know it won''t last
Somehow I know it won''t last too long

[Chorus]
She''ll come back to me
She''ll come back to me
She''ll come back to me
She''ll come back to me' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Italian Leather Sofa',  
'[Verse 1]
She doesn''t care whether or not he''s an island
She doesn''t care just as long as his ship''s coming in

[Chorus]
She doesn''t care whether or not he''s an island
They laugh they make money
He''s got a gold watch
She''s got a silk dress and healthy breasts
That bounce on his Italian leather sofa

[Verse 2]
She doesn''t care whether or not he''s a good man
She doesn''t care just as long as she still has her friends

[Chorus]
She doesn''t care whether or not he''s an island
They laugh, they make money
He''s got a gold watch
She''s got a silk dress and healthy breasts
That bounce on his Italian leather sofa

[Bridge]
She''s got a serrated edge
That she moves back and forth
It''s such a simple machine
She doesn''t have to use force
When she gets what she wants
She puts the rest on a tray in a ziplock bag

She''s got a serrated edge
That she moves back and forth
It''s such a simple machine
She doesn''t have to use force
When she gets what she wants
She puts the rest on a tray in a ziplock bag
...in the freezer

[Verse 1]
She doesn''t care whether or not he''s an island
She doesn''t care just as long as his ship''s coming in
Alright, here it comes, here it comes...

[Chorus]
She doesn''t care whether or not he''s an island
They laugh, they make money
He''s got a gold watch
She''s got a silk dress and healthy breasts
That bounce on his Italian leather sofa' 
);

INSERT INTO note (createdt, user_id, title, body) 
VALUES ('2020-02-01T14:00:00+08:00', 3, 
'Sad Songs & Waltzes',  
'I''m writing a song all about you
A true song as real as my tears
But you''ve no need to fear it
''Cause no one will hear it
Sad songs and waltzes
Aren''t selling this year

I''ll tell all about how you cheated
I''d like for the whole world to hear
I''d like to get even
With you ''cause you''re leavin''
But sad songs and waltzes
Aren''t selling this year

It''s a good thing that I''m not a star
You don''t know how lucky you are
Though my record may say it
No one will play it
Sad songs and waltzes
Aren''t selling this year

It''s a good thing that I''m not a star
You don''t know how lucky you are
Though my record may say it
No one will play it
Sad songs and waltzes
Aren''t selling this year' 
);

COMMIT;

