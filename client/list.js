Template.list.helpers({
   'items': function() {
      return Lists.find({
         'list': Session.get('selectedList'),
      },{
         sort: ['name']
      });
   }
});

Template.list.events({
   'click button.delete': function() {
      Lists.remove(this._id);
   }
});
