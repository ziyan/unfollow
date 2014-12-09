unfollow.namespace 'unfollow.network', (exports) ->
  'use strict'

  class Network
    constructor: (@network) ->
      @loading = []
      @expanding = []

      @index = {}
      @nodes = []
      @links = []
      @linked = {}

      @setup()

      if @network.data('id')
        @expand @network.data('id')
      else
        @bootstrap()

    setup: ->
      @svg = d3.select(@network[0]).append('svg')
      @defs = @svg.append('defs')
      @overlay = @svg.append('rect').attr('class', 'overlay')
      @g = [@svg.append('g'), @svg.append('g')]

      @link = @g[0].selectAll('.link')
      @node = @g[1].selectAll('.node')
      @pattern = @defs.selectAll('pattern')
      @zoom = d3.behavior.zoom().scaleExtent([1, 10]).on 'zoom', =>
        @g[0].attr('transform', 'translate(' + d3.event.translate + ')scale(' + d3.event.scale + ')')
        @g[1].attr('transform', 'translate(' + d3.event.translate + ')scale(' + d3.event.scale + ')')
      @svg.call(@zoom).on('dblclick.zoom', null)

      @force = d3.layout.force().charge(-120).linkDistance (data) ->
        return 15 + 15 * @nodes.length / 100 if data.bidirectional
        return 50 + @nodes.length / 2

      @force.on 'tick', =>
        @link
          .attr('x1', (data) -> data.source.x)
          .attr('y1', (data) -> data.source.y)
          .attr('x2', (data) -> data.target.x)
          .attr('y2', (data) -> data.target.y)
        @node
          .attr('cx', (data) -> data.x)
          .attr('cy', (data) -> data.y)
      @force.drag()
        .on 'dragstart', ->
          d3.event.sourceEvent.stopPropagation()
        .on 'dragend', (data) ->


      $(window).resize => @resize()
      @resize()

      # initial zoom
      initial_zoom = @network.data('zoom') or 4
      @zoom
        .translate([-@width / 2 * (initial_zoom - 1), -@height / 2 * (initial_zoom - 1)])
        .scale(initial_zoom)
        .event(d3.transition().duration(1000))

    resize: ->
      @width = $(window).width()
      @height = @network.data('height') or $(window).height()

      @svg.attr('width', @width).attr('height', @height)
      @overlay.attr('width', @width).attr('height', @height)
      @force.size([@width, @height])

      @update()

    update: ->
      @force.nodes(@nodes).links(@links).start()

      @pattern = @pattern.data(@nodes)
      @pattern.enter().append('pattern')
        .attr('id', (data) -> 'avatar_' + data.index)
        .attr('patternUnits', 'objectBoundingBox')
        .attr('width', 10)
        .attr('height', 10)
        .append('image')
        .attr('xlink:href', (data) -> data.node.avatar)
        .attr('x', 0)
        .attr('y', 0)
        .attr('width', 10)
        .attr('height', 10)
      @pattern.exit().remove()

      @node = @node.data(@nodes)
      @node.enter().append('circle')
        .attr('class', 'node')
        .attr('r', 5)
        .style('fill', (data) => 'url(#avatar_' + data.index + ')')
        .call(@force.drag)
        .on('dblclick', (data) =>
          return unless @network.data('id')
          data.fixed = true
          @expand data.node.id
        )
        .on('mouseover', (data) ->
        )
        .on('mouseout', (data) ->
        )
      @node.exit().remove()
      @node.append('title').text((data) -> data.node.name + ' (@' + data.node.screen_name + ')')
      @node
        .attr('class', (data) ->
          classes = ['node']
          classes.push 'expanded' if data.fixed
          classes.push 'verified' if data.node.verified
          classes.push 'protected' if data.node.protected
          classes.push 'discovered' if data.node.friends_ids or data.node.followers_ids
          classes.push 'loaded' if data.node.friends_loaded
          return classes.join(' ')
        )

      @link = @link.data(@links)
      @link.enter().append('line')
        .attr('class', 'link')
      @link.exit().remove()
      @link
        .attr('class', (data) ->
          classes = ['link']
          classes.push 'bidirectional' if data.bidirectional
          return classes.join(' ')
        )

    linker: ->
      linked = {}
      $.each @nodes, (index, data) =>
        return unless data.node.friends_ids
        for friend_id in data.node.friends_ids
          friend = @index[friend_id]
          continue unless friend

          bidirectional = false
          bidirectional = true if data.node.followers_ids and data.node.followers_ids.indexOf(friend.id) >= 0
          bidirectional = true if friend.friends_ids and friend.friends_ids.indexOf(data.node.id) >= 0

          link =
            source: data.node.index
            target: friend.index
            follower: data.node
            friend: friend
            bidirectional: bidirectional

          if bidirectional
            linked['' + data.node.id + '->' + friend.id] = null
            linked['' + friend.id + '->' + data.node.id] = null
            linked[[friend.id, data.node.id].sort().join('<->')] = link
          else
            key = '' + data.node.id + '->' + friend.id
            if linked[key] is undefined
              linked[key] = link 

      @links = []
      $.each linked, (i, link) =>
        @links.push link if link

    expand: (id) ->
      @expanding.push id
      @expander()

    expander: ->
      clearTimeout @expander_timer if @expander_timer
      @expander_timer = null

      expanding = []
      for id in @expanding
        node = @index[id]

        # check if node is now fully loaded
        if node
          continue if node.protected
          continue if node.friends_count is 0
          if node.friends_ids and node.friends_ids.length > 0
            loading = []
            for friend_id in node.friends_ids
              continue if @index[friend_id]
              loading.push friend_id
              break if loading.length > (@network.data('load') or 50)
            if loading.length > 0
              @load loading
            else
              node.friends_loaded = true
              @linker()
              @update()
            continue

        # continue to load node
        expanding.push id

        unfollow.ajax.post '/network/node', id: id, (data) =>
          return unless data and data.id and data.node
          node = data.node
          node.id = data.id

          existing = @index[node.id]
          if existing
            existing.friends_ids = node.friends_ids
            existing.followers_ids = node.followers_ids
            return

          @index[node.id] = node
          data = 
            node: node
            index: @nodes.length
          data.node.index = data.index
          @nodes.push data

          if data.index is 0
            data.px = @width / 2
            data.py = @height / 2
            data.fixed = true

          @linker()
          @update()

      @expanding = expanding
      @progress()

      return if @expanding.length is 0
      @expander_timer = setTimeout =>
        return if @network.parents('body').length is 0
        @expander()
      , 1000

    load: (ids) ->
      for id in ids
        @loading.push id
      @loader()

    loader: ->
      clearTimeout @loader_timer if @loader_timer
      @loader_timer = null

      loading = []
      for id in @loading
        # check if node is loaded
        continue if @index[id]

        # continue to load node
        loading.push id

      @loading = loading
      @progress()

      return if loading.length is 0

      unfollow.ajax.post '/network/nodes', ids: loading, (data) =>
        return unless data and data.nodes
        $.each data.nodes, (id, node) =>
          node.id = parseInt id

          existing = @index[node.id]
          if existing
            existing.friends_ids = node.friends_ids
            existing.followers_ids = node.followers_ids
            return

          @index[node.id] = node
          data = 
            node: node
            index: @nodes.length
          data.node.index = data.index
          @nodes.push data

          @linker()
          @update()

      @loader_timer = setTimeout =>
        return if @network.parents('body').length is 0
        @loader()
      , 1000

    progress: ->
      if @loading.length + @expanding.length > 0
        $('div.js-loading').addClass('progress-striped active')
      else
        $('div.js-loading').removeClass('progress-striped active')

    bootstrap: ->
      $.getJSON unfollow.settings.STATIC + '/' + unfollow.settings.VERSION + '/data/bootstrap.json', (data) =>
        @index = data

        $.each @index, (id, node) =>
          data = 
            node: node
            index: @nodes.length
          data.node.index = data.index
          if node.id is 7007262
            data.px = @width / 7 * 3
            data.py = @height / 2
            data.fixed = true
          if node.id is 1618521
            data.px = @width / 7 * 4
            data.py = @height / 2
            data.fixed = true
          @nodes.push data

        @linker()
        @update()

  on_network_init = ->
    network = $(this)
    network.data 'network', new Network(network)

  unfollow.init ->
    $('div.js-network').each on_network_init
