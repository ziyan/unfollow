DIR = ../app/locale
TARGETS = \
	$(DIR)/en_US/LC_MESSAGES/unfollow.mo \
	$(DIR)/zh_CN/LC_MESSAGES/unfollow.mo \
	$(DIR)/ja_JP/LC_MESSAGES/unfollow.mo \

.PHONY: all
all: $(TARGETS)

$(DIR)/%/LC_MESSAGES/unfollow.mo: %/LC_MESSAGES/unfollow.po
	@mkdir -p `dirname $@`
	@msgfmt -c -v -o $@ $<

%/LC_MESSAGES/unfollow.po: unfollow.pot
	@msgmerge $@ unfollow.pot > $@.tmp && mv $@.tmp $@

unfollow.pot: ../src ../html
	@rm -f $@
	@xgettext -d unfollow --force-po -o $@ --from-code utf-8 --language=perl -k -kGetText -s `find ../html -type f -iname \*.haml`
	@xgettext -d unfollow --force-po -o $@ --from-code utf-8 --language=c -k -kGetText -s -j `find ../src -type f -iname \*.go`

.PHONY: clean
clean:
	@rm -f unfollow.pot
	@rm -rf $(DIR)
