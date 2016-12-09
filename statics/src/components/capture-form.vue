<template>
  <div class="sub-left-panel">
    <div class="title-left-panel">New Capture :</div>
    <form @submit.prevent="start">
      <label>Name : </label></br>
      <input type="text" class="capture_input" v-model="name" /><br/>
      <label>Description : </label></br>
      <textarea type="text" class="capture_input" rows="2" v-model="desc"></textarea><br/>
      <label>Target : </label></br>
      <label class="radio-inline">
        <input type="radio" name="capture-target" value="selection" v-model="mode"> Nodes selection
      </label>
      <label class="radio-inline">
        <input type="radio" name="capture-target" value="gremlin" v-model="mode"> Gremlin Expression
      </label>
      <div v-if="mode == 'selection'">
        <node-selector v-model="node1" placeholder="Node 1"></node-selector>
        </br>
        <node-selector v-model="node2" placeholder="Node 2 (Optionnal)"></node-selector>
      </div>
      <br/>
      <div v-if="mode == 'gremlin'">
        <textarea type="text" class="capture_input" v-model="userQuery"></textarea><br/>
      </div>
      <button type="submit" class="btn btn-primary">Start</button>
      <button type="button" class="btn btn-danger" @click="reset">Cancel</button>
      <br/>
    </form>
  </div>
</template>

<script>
export default {

  data: function() {
    return {
      node1: "",
      node2: "",
      name: "",
      desc: "",
      userQuery: "",
      mode: "selection",
      queryError: "",
    };
  },

  components: {
    'node-selector': require('./node-selector.vue'),
  },

  computed: {

    query: function() {
      if (this.mode == "gremlin") {
        if (!this.userQuery) {
          this.queryError = "Gremlin query can't be empty";
          return;
        }
        return this.userQuery;
      } 
      else {
        if (!this.node1) {
          this.queryError = "At least one node has to be selected";
          return;
        }
        var q = "G.V().Has('TID', '" + this.node1 + "')";
        if (this.node2)
          q += ".ShortestPathTo(Metadata('TID', '" + this.node2 + "'), Metadata('RelationType', 'layer2'))";
        return q;
      }
    }

  },

  methods: {

    reset: function() {
      this.node1 = this.node2 = this.userQuery = "";
      this.name = this.desc = "";
    },

    start: function() {
      var self = this;
      if (!this.query) {
        $.notify({message: this.queryError}, {type: 'danger'});
        return;
      }
      $.ajax({
        dataType: "json",
        url: '/api/capture',
        data: JSON.stringify({"GremlinQuery": this.query, "Name": this.name, "Description": this.desc}),
        contentType: "application/json; charset=utf-8",
        method: 'POST',
        success: function() {
          $.notify({
            message: 'Capture created'
          },{
            type: 'success'
          });
          $("#capture").slideToggle(500, function () {});
          self.reset();
        },
        error: function(e) {
          $.notify({
            message: 'Capture create error: ' + e.responseText
          },{
            type: 'danger'
          });
        }
      });
    },

  }

}
</script>
