import re
import json
from PyPDF2 import PdfReader

# Just a comment for testing render deploy on push

def convert_to_num(text_num):
    try:
        text_list = text_num.split(".")
        return 10000 * int(text_list[0]) + 1000 * int(text_list[1]) + int(text_list[2])
    except:
        return 0

reader = PdfReader("./admin/Savchenko_2008.pdf")

to_find = "\d{1,2}\.\d{1,3}\.\d{1,3}"
i = 0
numbers = set()
for page in reader.pages:
    i += 1
    text = page.extract_text() 
    res_search = re.findall(to_find, text)
    if res_search:
        for j in res_search:
            numbers.add(j)

numbers = sorted(list(numbers), key = lambda x: convert_to_num(x))[:-1]

result_dict = {}

for item in numbers:
    parts = item.split('.')
    major, minor, revision = map(int, parts)
    
    if major not in result_dict:
        result_dict[major] = {}
    
    if minor not in result_dict[major]:
        result_dict[major][minor] = []
    
    result_dict[major][minor].append(item)


with open('problems_list.json', 'w') as f:
    json.dump(result_dict, f)
