Template.add.events({
   'submit': function(e,t) {
      e.preventDefault();
      e.stopPropagation();
      var val = s.trim(t.find('input[name=site]').value);
      if( val === ''  ) { return false; }

      if( Lists.find({
         list: Session.get('selectedList'),
         name: val
      }).count() > 0 ) { return false; }

      Lists.insert({
         list: Session.get('selectedList'),
         name: val
      }, function(err) {
         if(err) {
            console.error(err);
            return false;
         }
         e.currentTarget.reset();
      });
      return false;
   },
});
