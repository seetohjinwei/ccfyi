VENV=../venv
REQUIREMENTS=../requirements.txt
PYTHON=$(VENV)/bin/python3
PIP=$(VENV)/bin/pip

.PHONY: all

all: run

run: $(VENV)/bin/activate
	$(PYTHON) main.py

test: $(VENV)/bin/activate
	$(PYTHON) -m unittest discover tests

$(VENV)/bin/activate: $(REQUIREMENTS)
	python3 -m venv $(VENV)
	$(PIP) install -r $(REQUIREMENTS)

clean:
	rm -rf __pycache__
	rm -rf $(VENV)
