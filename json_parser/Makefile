VENV=../venv
REQUIREMENTS=../requirements.txt
PYTHON=$(VENV)/bin/python3
PIP=$(VENV)/bin/pip
EDITOR=nvim
VENV_ACTIVATE=. $(VENV)/bin/activate

.ONESHELL:
.PHONY: all

all: edit

edit: $(VENV)/bin/activate
	$(VENV_ACTIVATE) && $(EDITOR)

run: $(VENV)/bin/activate
	$(VENV_ACTIVATE) && $(PYTHON) main.py $(input)

test: $(VENV)/bin/activate
	if [ -z "$(test)" ]; then \
		$(VENV_ACTIVATE) && $(PYTHON) -m unittest discover tests; \
	else \
		$(VENV_ACTIVATE) && $(PYTHON) -m unittest $(test); \
	fi

$(VENV)/bin/activate: $(REQUIREMENTS)
	rm -rf $(VENV)
	python3 -m venv $(VENV)
	$(VENV_ACTIVATE) && $(PIP) install -r $(REQUIREMENTS)

clean:
	rm -rf __pycache__
	rm -rf $(VENV)
