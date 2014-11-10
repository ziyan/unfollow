unfollow.namespace 'unfollow.utils', (exports) ->
  'use strict'

  # setup google analytics
  analytics = (account) ->
    return unless account

    require ['ga'], ->
      window.ga 'create', account, 'unfollow.com'
      window.ga 'send', 'pageview',
        page: location.pathname + location.search + location.hash
        location: location.href

  # extract data attributes on an element into an object
  exports.data = (element, filter) ->
    filter = filter or -> false
    data = {}
    $.each element[0].attributes, (index, attribute) ->
      return unless /^data-/.test(attribute.name)
      return if filter(attribute.name)
      data[attribute.name.substr(5).replace(/\-/g, '_')] = attribute.value
    return data

  unfollow.init ->
    analytics(unfollow.settings.ANALYTICS)
