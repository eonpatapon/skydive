var FilterButton = {

  template: '<div class="btn-group" role="group">\
              <button class="btn btn-default btn-sm" readonly>\
                {{ label }} : {{ value }}\
              </button>\
              <button class="btn btn-danger btn-sm" v-on:click="remove">\
                <span class="glyphicon glyphicon-remove"></span>\
              </button>\
             </div>',

  props: ['label', 'value'],

  methods: {

    remove: function() {
      this.$emit('remove', this.label, this.value);
    }

  }

};

Vue.component('filter-list', {

  template: '\
      <div class="filter">\
        <div class="filter-form">\
          <form class="form-inline" @submit.prevent="add(label, value)">\
            <div class="form-group">\
              <autocomplete :suggestions="suggestedKeys"\
                            :placeholder="labelPlaceholder"\
                            v-model="label"\
              ></autocomplete>\
            </div>\
            <div class="form-group">\
              <autocomplete :suggestions="suggestedValues"\
                            :placeholder="valuePlaceholder"\
                            v-model="value"\
                            @select="add(label, value)"\
              ></autocomplete>\
            </div>\
            <button class="btn btn-sm btn-primary" type="submit">\
              <span class="glyphicon glyphicon-plus"></span>\
            </button>\
          </form>\
        </div>\
        <div class="filter-list">\
          <filter-button v-for="filter in filterList"\
                         :label="filter.label"\
                         :value="filter.value"\
                         @remove="remove">\
          </filter-button>\
        </div>\
      </div>',

  components: {
    'filter-button': FilterButton,
  },

  data: function() {
    return {
      filters: {},
      value: "",
      label: "",
    };
  },

  props: {
    keySuggestions: {
      type: Function
    },
    valueSuggestions: {
      type: Function
    },
    labelPlaceholder: {
      type: String
    },
    valuePlaceholder: {
      type: String
    }
  },

  computed: {

    filterList: function() {
      var self = this;
      return Object.keys(this.filters).reduce(function(acc, l) {
        for (var v of self.filters[l]) {
          acc.push({label: l, value: v});
        }
        return acc;
      }, []);
    },

    suggestedKeys: function() {
      return this.keySuggestions || [];
    },

    suggestedValues: function() {
      var self = this;
      if (this.valueSuggestions) {
        return function() {
          if (self.label) {
            return self.valueSuggestions(self.label);
          }
          else {
            var d = $.Deferred();
            d.resolve([]);
            return d;
          }
        };
      }
      else {
        return [];
      }
    },

  },

  methods: {

    add: function(label, value, notify) {
      if (this.filters[label] && value in this.filters[label])
        return;
      if (!this.filters[label])
        Vue.set(this.filters, label, [value]);
      else
        this.filters[label].push(value);

      // reset form values
      this.value = this.label = "";
      
      if (typeof notify == "undefined" || notify === true) {
        this.$emit('update', this.filters);
      }
    },

    remove: function(label, value) {
      if (!this.filters[label])
        return;
      for (var i in this.filters[label]) {
        if (this.filters[label][i] == value) {
          this.filters[label].splice(i, 1);
          if (this.filters[label].length === 0) {
            Vue.delete(this.filters, label);
          }
          this.$emit('update', this.filters);
          break;
        }
      }
    }

  }

});
