FROM python:3.9

WORKDIR /webapp

COPY requirements.txt /webapp/requirements.txt

RUN pip install --no-cache-dir --upgrade -r /webapp/requirements.txt

COPY . .

CMD [ "python", "app.py" ]