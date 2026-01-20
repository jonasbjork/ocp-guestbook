2026-01-15
    Idag har jag sutit och försökt att få igång alla de fya Containers som frontend, backend, Postgresql och Redis. Men det tog mig tid och jag klarade att bygga en container för Frontend och testade att bygga Backend i terminalen bara att för att se om jag får "healthy" för min postgresql och redis. 

    Frontend gick bra men när det kom till backend, postgresql och redis var det svårt  få de kommunicera. 
    Problemet var att jag hade inte skapat nätverk och det var självklarhet att det inte gick men sen när jag väl lärde mig, tog jag bort backend containern och och byggde om det men "--network" flaga

2026-01-16
    Dagens mål var att jag och Jonas ska kunna skapa container filer för Backend, Postgresql och Redis och bygga de i vårat repo. Det blev mycket enkklare när jag hade redan testat det i terminalen dagen innan och visste hur byggar man till exampel en postgresql container och vilka variablar behöver man där. 

    Men det som jag stötte på och tog längre tid var om hur ska jag bygga nätverket för de i koden. I terminalen så skrev jag podman create --network guestbook-net och det gick fort. Efter en stund undersökningar så fick jag reda på att när jag är klar med skriva de tre Conatinerfiles då kan jag bara lägga till nätverk i det kommandot som jag startar alla mina Conatiner med. Detta gör att jag kan ha alla fyra Conatiner under en och samma nätverk.
    Exampel på Comandot "podman run -d --name Contanier-name --network nework-name -p port Container-Name"

2026-01-17
    Idag har jag och Jonas sutit och samtalt alla och försökte få igång application och kunna skriva i frontend och skicka den vidare. samt att den ska sparas så länge continerna inte är borttagen. 

    Problemet idag var att vi kunde komma till frontend och backend i terminalen fick vi healthy med hjälp av curl men när vi skrev i hemsidan så kunde den inte skicka eller spara. 

    Vi har undersökt i flera timmar och kollat i google och AI och till slut fick vi reda på att i main.go bevhöver vi specificera databasen med postgresql och även för redis behöver vi specificera redis, port och lösenord. Sedan sätte vi resolver proxy i nginx.conf till localhost och testade igen.

    Nu har vi igång alla container och allt sparas så länge man inte ta bort Containers från podman. 


2026-01-19
    Jag och Jonas har sutit idag och dokumenterad samt fyllde i logbooken. Vi har diskuterat om hur hade vi kunnat lösa det på ett annat sätt, om det finns ett annat sätt och hur hade det påverkat produktionen. Vi har några ideér men vi vill gärna diskutera det imorgon när vi har lektion och träffas på plats med Jonas Björk (Läraren). 

    
