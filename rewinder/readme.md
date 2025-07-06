As a novice to the android developer experience, I am learning as I go.

This is a non-exhaustive list of issues that I have encountered. 

# Java mismatch. 

I upgraded java on my machine, and it became immediately incompatible with the gradle version for the project.
The errors were all a mess with things such as - 

"""BUG! exception in phase 'semantic analysis' in source unit '_BuildScript_' Unsupported class file major version 65"""

Fix:

I managed to fix this by switching back to a previous
java version by running:
        
        # print out java versions
        update-java-alternatives --list
        # set a lower java version
        sudo update-java-alternatives --set // version 17 or whatever

After a few more updates here and there I finally managed to get something that built it normally.
I think that the issue is that everything besides the java version is versioned. I wonder if there
is a way to also version the 

Long term solution:

I added the acceptable version of java to /thirdparty, and created a docker target for building the app.

