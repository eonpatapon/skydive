var VueTopology = {
  el: "#vue-topology",

  data: {
    defaultTopologyFilters: [
      ['State', 'UP']
    ]
  },

  // init default filters
  mounted: function() {
    for (var f of this.defaultTopologyFilters) {
      this.$refs.topologyFilters.add(f[0], f[1], false);
    }
  },

  methods: {

    topologyFilters: function() {
      // access component via ref
      // available once the vue is mounted
      if (this.$refs && this.$refs.topologyFilters)
        return this.$refs.topologyFilters.filters;
      else
        return {};
    },

    topologyFiltersKeys: function() {
      var query = "G.V().Keys().Dedup()";
      return $.ajax({
        dataType: "json",
        url: '/api/topology',
        data: JSON.stringify({"GremlinQuery": query}),
        method: 'POST',
        error: function(e) {
          $.notify({
            message: 'Gremlin request error: ' + e.responseText
          },{
            type: 'danger'
          });
        }
      });
    },

    topologyFiltersValues: function(key) {
      var query = "G.V().Values('"+key+"').Dedup()";
      return $.ajax({
        dataType: "json",
        url: '/api/topology',
        data: JSON.stringify({"GremlinQuery": query}),
        method: 'POST',
        error: function(e) {
          $.notify({
            message: 'Gremlin request error: ' + e.responseText
          },{
            type: 'danger'
          });
        }
      }).then(function(data) {
        return data.map(function(v) {
          return v.toString();
        });
      });
    },

    applyTopologyFilters: function() {
      var self = this;
      topologyLayout.nodes.map(function(n) {
        self.applyTopologyFiltersNode(n);
      });
      topologyLayout.Redraw();
    },

    applyTopologyFiltersNode: function(node) {
      var filters = this.topologyFilters();
      if (Object.keys(filters).length === 0) {
        node.Visible = true;
        node.Filtered = false;
        return node;
      }

      var match = [];
      for (var filter in filters) {
        if (!node.Metadata[filter]) {
          match.push(null);
        } else {
          match.push(filters[filter].indexOf(node.Metadata[filter]) !== -1);
        }
      }
      var visible = match.every(function(e) { return e === true || e === null; });
      var filtered = match.every(function(e) { return e === true; });

      node.Visible = visible;
      node.Filtered = filtered;
      return node;
    }

  }

};
