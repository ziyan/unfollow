unfollow.namespace 'unfollow.locale', (exports) ->
  'use strict'

  translations =
    'zh_CN':
      'Test': '测试'
      'Users': '用户'
      'Networks': '网络'
      'Search for': '搜索'
      '%d hours ago': '%d小时前'
      'An hour ago': '一小时前'
      '%d minutes ago': '%d分钟前'
      'Just now': '刚才'
      'Few minutes ago': '几分钟前'
      'Today at %s': '今天 %s'
      'Yesterday': '昨天'
      '%Y/%m/%d': '%Y年%m月%d日'
      '%m/%d': '%m月%d日'

  on_locale_click = (e) ->
    locale = $(this).data 'locale'
    return if locale == unfollow.settings.LOCALE
    $.cookie 'locale', locale,
      expires: 365
    window.location.replace(window.location.href)

  exports.gettext = (text) ->
    locale = unfollow.settings.LOCALE
    return text if locale not of translations
    return text if text not of translations[locale]
    return translations[locale][text]

  unfollow.init ->
    $('a.js-locale').on 'click', on_locale_click

$ ->
  window.gettext = unfollow.locale.gettext
  debug.debug 'unfollow.locale.init', gettext 'Test'

