/* jshint multistr: true */

Vue.component('inject-form', {

  mixins: [apiMixin, notificationMixin],

  template: '\
    <form @submit.prevent="inject">\
      <div class="form-group">\
        <label for="inject-type">Type</label>\
        <select id="inject-type" v-model="type" class="form-control input-sm">\
          <option value="icmp">ICMPv4/Echo request</option>\
          <option value="tcp">TCP ping</option>\
        </select>\
      </div>\
      <div class="form-group" v-if="type == \'tcp\'">\
        <label for="inject-src-port">Src Port</label>\
        <input id="inject-src-port" type="number" class="form-control input-sm" v-model="srcPort" min="1" />\
      </div>\
      <div class="form-group" v-if="type == \'tcp\'">\
        <label for="inject-dst-port">Dst Port</label>\
        <input id="inject-dst-port" type="number" class="form-control input-sm" v-model="dstPort" min="1" />\
      </div>\
      <div class="form-group">\
        <label>Target</label>\
        <node-selector class="inject-target"\
                       placeholder="From"\
                       attr="ID"\
                       v-model="node1"></node-selector>\
        <node-selector placeholder="To"\
                       attr="ID"\
                       v-model="node2"></node-selector>\
      </div>\
      <div class="form-group">\
        <label for="inject-count">Nb. of packets</label>\
        <input id="inject-count" type="number" class="form-control input-sm" v-model="count" min="1" />\
      </div>\
      <button type="submit" class="btn btn-primary">Inject</button>\
      <button type="button" class="btn btn-danger" @click="reset">Reset</button>\
    </form>\
  ',

  data: function() {
    return {
      node1: "",
      node2: "",
      count: 1,
      type: "icmp",
      srcPort: 0,
      dstPort: 0,
    };
  },

  created: function() {
    if (this.$store.state.currentNode) {
      this.node1 = this.$store.state.currentNode.ID;
    }
  },

  beforeDestroy: function() {
    // FIXME: we should just call reset() here,
    // but the watchers are not being evaluated :/
    if (this.node1) {
      this.highlightNode(this.node1, false);
    }
    if (this.node2) {
      this.highlightNode(this.node2, false);
    }
  },

  computed: {

    error: function() {
      if (!this.node1 || !this.node2) {
          return "Source and destination interfaces must be selected";
      } else {
          return;
      }
    },

  },

  watch: {

    node1: function(newVal, oldVal) {
      if (oldVal) {
        this.highlightNode(oldVal, false);
      }
      this.highlightNode(newVal, true);
    },

    node2: function(newVal, oldVal) {
      if (oldVal) {
        this.highlightNode(oldVal, false);
      }
      this.highlightNode(newVal, true);
    }

  },

  methods: {

    highlightNode: function(id, bool) {
      if (bool)
        this.$store.commit('highlight', id);
      else
        this.$store.commit('unhighlight', id);
    },

    reset: function() {
      var self = this;
      this.node1 = this.node2 = "";
      this.count = 1;
      this.type = "icmp";
    },

    inject: function() {
      var self = this;
      if (this.error) {
        this.$error({message: this.error});
        return;
      }
      $.ajax({
        dataType: "json",
        url: '/api/injectpacket',
        data: JSON.stringify({
          "Src": "G.V('" + this.node1 + "')",
          "Dst": "G.V('" + this.node2 + "')",
          "SrcPort": this.srcPort,
          "DstPort": this.dstPort,
          "Type": this.type,
          "Count": this.count
        }),
        contentType: "application/json; charset=utf-8",
        method: 'POST',
      })
      .then(function() {
        self.$success({message: 'Packet injected'});
      })
      .fail(function(e) {
        self.$error({message: 'Packet injection error: ' + e.responseText});
      });
    },

  }

});
