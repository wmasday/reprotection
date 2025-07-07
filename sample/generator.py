import os
import random
import string

blocked_keywords = [
    "porno", "porn", "bokep", "sex", "seks", "18+", "nude", "nudity", "hot",
    "jav", "anal", "bdsm", "payudara", "pantat", "bugil", "mesum", "masturbasi",
    "bunuh", "pembunuhan", "tembak", "darah", "mati", "pembantaian", "tusuk",
    "bakar", "bom", "teroris", "ledakan", "pukulan", "gantung diri",
    "kafir", "babi", "anjing", "goblok", "bodoh", "jancok", "tolol", "kontol",
    "asu", "bangsat", "idiot", "penghinaan", "rasis", "fasis", "nazi",
    "gratis", "free", "diskon", "promo", "klik di sini", "penipuan", "bitcoin gratis",
    "investasi cepat", "cepat kaya", "pinjaman online", "uang cepat", "hacks", "cheat",
    "memek", "ngentot", "coli", "kampret", "taik", "tai", "peler", "titit",
    "setan", "iblis", "anjrit", "keparat"
]

output_dir = "output_html"
os.makedirs(output_dir, exist_ok=True)
random_words = ['gacor', 'elite', 'pro', 'xd', 'v2', 'main', 'demo', 'lol', 'dark', 'lite']

for i in range(100):
    keyword = random.choice(blocked_keywords).replace(" ", "-").replace("+", "plus")
    extra = random.choice(random_words)
    number = ''.join(random.choices(string.digits, k=3))
    filename = f"{keyword}-{extra}-{number}.html"
    filepath = os.path.join(output_dir, filename)

    with open(filepath, "w") as f:
        f.write(f"<html><body><h1>{filename}</h1></body></html>")
