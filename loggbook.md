## 2026-01-15
    Idag har jag sutit och försökt att få igång alla de fya Containers som frontend, backend, Postgresql och Redis. Men det tog mig tid och jag klarade att bygga en container för Frontend och testade att bygga Backend i terminalen bara att för att se om jag får "healthy" för min postgresql och redis. 

    Frontend gick bra men när det kom till backend, postgresql och redis var det svårt  få de kommunicera. 
    Problemet var att jag hade inte skapat nätverk och det var självklarhet att det inte gick men sen när jag väl lärde mig, tog jag bort backend containern och och byggde om det men "--network" flaga

## 2026-01-16
    Dagens mål var att jag och Jonas ska kunna skapa container filer för Backend, Postgresql och Redis och bygga de i vårat repo. Det blev mycket enkklare när jag hade redan testat det i terminalen dagen innan och visste hur byggar man till exampel en postgresql container och vilka variablar behöver man där. 

    Men det som jag stötte på och tog längre tid var om hur ska jag bygga nätverket för de i koden. I terminalen så skrev jag podman create --network guestbook-net och det gick fort. Efter en stund undersökningar så fick jag reda på att när jag är klar med skriva de tre Conatinerfiles då kan jag bara lägga till nätverk i det kommandot som jag startar alla mina Conatiner med. Detta gör att jag kan ha alla fyra Conatiner under en och samma nätverk.
    Exampel på Comandot "podman run -d --name Contanier-name --network nework-name -p port Container-Name"

## 2026-01-17
    Idag har jag och Jonas sutit och samtalt alla och försökte få igång application och kunna skriva i frontend och skicka den vidare. samt att den ska sparas så länge continerna inte är borttagen. 

    Problemet idag var att vi kunde komma till frontend och backend i terminalen fick vi healthy med hjälp av curl men när vi skrev i hemsidan så kunde den inte skicka eller spara. 

    Vi har undersökt i flera timmar och kollat i google och AI och till slut fick vi reda på att i main.go bevhöver vi specificera databasen med postgresql och även för redis behöver vi specificera redis, port och lösenord. Sedan sätte vi resolver proxy i nginx.conf till localhost och testade igen.

    Nu har vi igång alla container och allt sparas så länge man inte ta bort Containers från podman. 


## 2026-01-19
    Jag och Jonas har sutit idag och dokumenterad samt fyllde i logbooken. Vi har diskuterat om hur hade vi kunnat lösa det på ett annat sätt, om det finns ett annat sätt och hur hade det påverkat produktionen. Vi har några ideér men vi vill gärna diskutera det imorgon när vi har lektion och träffas på plats med Jonas Björk (Läraren). 

## 2026-01-20
    Tisdag idag och under lektionen fick vi att jobba med skapa workflows för frontend och backend. Jag och Jonas tänkte att bygga båda podman och docker baserad workflows för att se om det kommer funka eller inte. 
    Vi började med frontend först och Jonas valde docker och jag började med att bygga frontend workflowet i podman för att se hur det reagerar och om det ens gå. 

    Jag tänkte att eftersom vi har jobbat med podman sen början så kan det vara bra att fortsätta med podman även i workflows om det går. 
     - Jag behövde ha ett steps för att isntallera och Podman allra först efter att den har     gjort checkout. 
     - logga in på Github container registry med Podman
     - Bygga imagen
     - Pusha imagen

    När jag har kommitat ändringar och pushat märkte vi att workflowet gick inte genom och under loggar stod tydlig att mitt ´username´ på Github börja med storbokstav (Narah111).
    Jag har aldrig haft detta problem men jag började googla och se om jag hittar en lösning och efter stund tog jag och Jonas en dialog med Läraren. Till slut har vi hårdkodade usernmet i image prefix med lite bokstav och då gick min workflows genom. 

    Vi har lärt oss något nytt i den delen och även provat att frontenden går och köra på podman.


    Vi började lite med Backend men då var dagen slut och vi bestämde att fortsatta imorgon och vi ska träffas på skolan själva.

## 2026-01-21
    Dagen mål är att jag och Jonas ska sätta upp även backend workflowet så att vi är ikapp tills torsdagen lektion. 

    Vi började med samma plan att Jonas ska bygga backend med att docker och jag ska bygga den med podman. Vi tänkte prova att skriva samma workflows som frontend fast ändra vissa värden som build bakcend image osv bara för att se om det funkar. 

    Det funkade tyvärr inte vilket är inte så konstigt och vi förväntade oss det.
    Efter en stund google och AI fråga stund fick vi redan på att vi behöver Qemu setup sedan definera platforms under build och push steps då kunde vi bygga imagen i både Linux och Mac miljö(amd64 & arm64) och testade workflowes och det gick. Fast det 8-9min. 

    Efter en stund diskution, kom Jonas med förslaget att vi får rita på tavlan hela flödet från början till slut för att se hur vi tänker och hur kan vi undvika den långa vänte tiden. 

![alt text](F15AC216-F7B2-4F7A-8118-078C41469845_1_105_c-1.jpeg)

    Som bilden visar så har vi tänkt att vi ska ha en future branch (development) och bygga backend workflows att när det är push till dev ska den bygga bara arm64 och när det är PR till main ska den bygga amd64 för så att vi kan använda den i Openshift. Vi tror att med hjälp av den här fördelningen kan vi spara både tid och gör det mer specificerad. 

    Juste, jag började bygga backenden med Docker då fick jag inga bra rekomendationer i framtiden varken från läraren eller i google. 

    Jag och Jonas har skrivit om backenden i två olika sätt men samma funktion och efter ett par timmars slit lyckades vi ha en fungerande backend som vi önkade enligt bilden.

## 2026-01-22 - 2026-01-23 2026-01-26
    Den 22 januri hade vi lektion och denna lektionen gick läraren genom Openshift generellt och de viktiga funktioner som vi kommer behöver använda i vårat prjekt. 
    Den 23de(fredag) och 26te(måndag) fick vi jobba med labben i RedHat som vi kan öva på den en hel del. 

## 2026-01-27
    Lektion men jag tyvärr har missat det pga av sjukdom. 

