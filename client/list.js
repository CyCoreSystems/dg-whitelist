Template.list.helpers({
   'items': function() {
      return Lists.find({
         'list': Session.get('selectedList'),
      },{
         sort: ['name']
      });
   },
   'tableSettings': function() {
      return {
         fields: [
            {
               label: 'Site',
               key: 'name',
            },
            {
               label: 'Added',
               key: 'added',
               fn: function(val) {
                  return moment.unix(val).calendar();
               },
               sortByValue: true
            },
            {
               label: 'Delete',
               key: '_id',
               fn: function(val) {
                  return new Spacebars.SafeString('<button class="btn delete" data-id='+val+'><i class="fa fa-trash"></i></button>');
               },
               sortable: false
            },
         ]
      };
   }
});

Template.list.events({
   'click .delete': function(e) {
      e.preventDefault();
      e.stopPropagation();
      Lists.remove(e.currentTarget.dataset.id);
   }
});
