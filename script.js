if (Array.isArray(input.network)) {
  var matches = input.network.filter(function (line) {
    if (line.left === "BILLED_TO" && line.right === "BILLED_TO" && line.link === "CreditCard" && line.overusers > 0) {
      return true;
    }
  });
  matches.length > 0;
}