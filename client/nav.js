Session.setDefault('selectedList','whitelist');

Template.nav.helpers({
   'listTypes': function() {
      return ListTypes.find();
   },
   'isActive': function(listType) {
      console.log('This listType:',listType);
      if(Session.equals('selectedList',listType)) {
         console.log("match");
         return 'active';
      }
      return '';
   }
});

Template.nav.events({
   'click li': function(e,t) {
      e.stopPropagation();
      e.preventDefault();
      Session.set('selectedList',this.name);
   }
});
