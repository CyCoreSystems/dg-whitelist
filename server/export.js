function exportList(listType) {
   dg.export(listType,Lists.find({list: listType}));
}
