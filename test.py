#req.text

text = '[{"id":4,"name":"Andrei"},{"id":5,"name":"Cosmin"},{"id":6,"name":"Teodor"},{"id":7,"name":"Costelusu"},{"id":8,"name":"lelelelele"},{"id":9,"name":"bossu1"},{"id":10,"name":"bossu123"},{"id":11,"name":"bossu1234"},{"id":12,"name":"coiosu"}]'

text = text.split(',')

for x in range(len(text)):
    if x%2==0:
        print (text[x].split(':')[1])
        print (text[x+1].split(':')[1].split('"')[1])
        print ("+++++++++++++++++")
