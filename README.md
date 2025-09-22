# Distribute


> [!WARNING]
> I am VERY new to GO Language, this is literally my first project with this language and I will
> keep this warning here, till I feel confident that this project has no **amateur** bugs. Please
> run this script ONLY on not so important systems. For e.g. I run it to connect my personal
> Raspberry pi 5 which is not connected to external internet. 


# Problem statement

I love developing websites and apps and one of the main problems I faced was to test these apps
and websites on my phone. I did not want to open my laptop's ports, so I got a raspberry pi to
setup a home server. Now, the trouble was that I had manually SFTP into my pi and then transfer
the build files and everything. It was tedious and repetitive, and don't ask me how frustrated I
got when I had to check small small changes on mobile devices, the amount of times I had to build
and transfer.


# The Solution

I initially built a bash script, with the help of chatgpt, that did all of this, but it was very
bad in terms of presentation and the syntax didnt make sense to me. So, I decided to rewrite it on
my own, in a language that I know! But wait! I had found [gum by
charmbracelet](https://github.com/charmbracelet/gum) That was beautiful and that led me to [Bubble
Tea](https://github.com/charmbracelet/bubbletea) and I thought to myself why not build it using Go
and this way I would also, learn about the language. So, that's how this project was born
