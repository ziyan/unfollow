DIR = ../app/static/css
UNFOLLOW_TARGET = unfollow.css
UNFOLLOW_SOURCES = \
	bootstrap/bootstrap.css \
	fontawesome/fontawesome.css \
	unfollow/unfollow.sass.css \

.PHONY: all
all: debug

.PHONY: debug
debug: $(DIR)/$(UNFOLLOW_TARGET)

.PHONY: release
release: $(DIR)/$(UNFOLLOW_TARGET)
	@cat $(UNFOLLOW_TARGET) | sed 's/\/\*!/\/**/g' | yui-compressor --type css > $(DIR)/$(UNFOLLOW_TARGET)

$(DIR)/$(UNFOLLOW_TARGET): $(UNFOLLOW_TARGET)
	@mkdir -p $(DIR)
	@cat $(UNFOLLOW_TARGET) > $(DIR)/$(UNFOLLOW_TARGET)

$(UNFOLLOW_TARGET): $(UNFOLLOW_SOURCES)
	@cat $(UNFOLLOW_SOURCES) > $(UNFOLLOW_TARGET)

%.sass.css: %.sass $(shell find . -type f -iname _\*.sass)
	@cat $< | sass --compass --load-path . > $<.css

.PHONY: clean
clean:
	@find . -iname \*.sass.css -exec rm -f {} \;
	@rm -rf .sass-cache
	@rm -f $(UNFOLLOW_TARGET)
	@rm -rf $(DIR)

