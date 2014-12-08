DIR = ../app/static/js
UNFOLLOW_TARGET = unfollow.js
UNFOLLOW_SOURCES = \
	debug.js \
	jquery/jquery.js \
	jquery/jquery.cookie.coffee.js \
	jquery/jquery.livequery.js \
	jquery/jquery.imagesloaded.js \
	jquery/jquery.pjax.coffee.js \
	bootstrap/bootstrap.js \
	d3.js \
	require.js \
	unfollow/unfollow.coffee.js \
	unfollow/unfollow.locale.coffee.js \
	unfollow/unfollow.utils.coffee.js \
	unfollow/unfollow.ajax.coffee.js \
	unfollow/unfollow.network.coffee.js \

.PHONY: all
all: debug

.PHONY: debug
debug: $(DIR)/$(UNFOLLOW_TARGET)

.PHONY: release
release: $(DIR)/$(UNFOLLOW_TARGET)
	@cat $(UNFOLLOW_TARGET) | sed 's/\/\*!/\/**/g' | yui-compressor --type js > $(DIR)/$(UNFOLLOW_TARGET)

$(DIR)/$(UNFOLLOW_TARGET): $(UNFOLLOW_TARGET)
	@mkdir -p $(DIR)
	@cat $(UNFOLLOW_TARGET) > $(DIR)/$(UNFOLLOW_TARGET)

$(UNFOLLOW_TARGET): $(UNFOLLOW_SOURCES)
	@cat $(UNFOLLOW_SOURCES) > $(UNFOLLOW_TARGET)

%.coffee.js: %.coffee
	@cat $< | coffee -c -s > $<.js

.PHONY: clean
clean:
	@find . -iname \*.coffee.js -exec rm -f {} \;
	@rm -f $(UNFOLLOW_TARGET)
	@rm -rf $(DIR)
