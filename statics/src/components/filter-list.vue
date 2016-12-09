<template>
  <div class="filter">
        <div class="filter-form">
          <form class="form-inline" @submit.prevent="add(label, value)">
            <div class="form-group">
              <autocomplete :suggestions="suggestedKeys"
                            :placeholder="labelPlaceholder"
                            v-model="label"
              ></autocomplete>
            </div>
            <div class="form-group">
              <autocomplete :suggestions="suggestedValues"
                            :placeholder="valuePlaceholder"
                            v-model="value"
                            @select="add(label, value)"
              ></autocomplete>
            </div>
            <button class="btn btn-sm btn-primary" type="submit">
              <span class="glyphicon glyphicon-plus"></span>
            </button>
          </form>
        </div>
        <div class="filter-list">
          <filter-button v-for="filter in filterList"
                         :label="filter.label"
                         :value="filter.value"
                         @remove="remove">
          </filter-button>
        </div>
      </div>
</template>

<script>
export default {

  components: {
    'filter-button': require('./filter-button.vue'),
  },

  data() {
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

    filterList() {
      return Object.keys(this.filters).reduce((acc, l) => {
        for (var v of this.filters[l]) {
          acc.push({label: l, value: v});
        }
        return acc;
      }, []);
    },

    suggestedKeys() {
      return this.keySuggestions || [];
    },

    suggestedValues() {
      if (this.valueSuggestions) {
        return () => {
          if (this.label) {
            return this.valueSuggestions(this.label);
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

    add(label, value, notify) {
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

    remove(label, value) {
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

}

</script>
