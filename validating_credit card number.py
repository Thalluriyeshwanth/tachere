import re

def validate_credit_card(card_number):
    pattern = r'^[456]\d{3}(-?\d{4}){3}$'
    if re.match(pattern, card_number):
        card_number = card_number.replace('-', '')
        if re.search(r'(\d)\1{3}', card_number):
            return 'Invalid'
        else:
            return 'Valid'
    else:
        return 'Invalid'

if __name__ == "__main__":
    n = int(input())
    for _ in range(n):
        card = input().strip()
        print(validate_credit_card(card))
