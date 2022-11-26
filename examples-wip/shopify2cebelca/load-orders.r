print {
      name: "script to load shopify orders"
}

print order-url: https://37d36ed5bd8946677928db24cecc062d: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@cebelcabiztest.myshopify.com/admin/api/2020-10/orders.json

data: print get order-url

drop-in "after loading webpage" { }

// TODO -- make json load the data into raw-map (dicts) and list (arrays) or blocks and contexts ... think about it or try the fastest one for now
// TODO -- make a validation so that it ignores the keys that aren't included in it
// TODO -- make function infer-validation that accepts a rawmap and returns a validation rules
// TODO LATER -- how do we retrieve or validate on children .. for now we probably just get the node out manually and validate it flat
// TODO LATER -- how to validate an array of dicts ... probably best with just a map for now ... maybe map-check or something for error handling

