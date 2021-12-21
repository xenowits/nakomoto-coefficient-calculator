const range = (len) => {
  const arr = [];
  for (let i = 0; i < len; i++) {
    arr.push(i);
  }
  return arr;
};

const data = [
  {
    chainName: "Binance",
    tokenName: "BNB",
    currVal: 7,
    prevVal: 7,
    changeVal: 0,
  },
  {
    chainName: "Cosmos",
    tokenName: "ATOM",
    currVal: 7,
    prevVal: 7,
    changeVal: 0,
  },
  {
    chainName: "Osmosis",
    tokenName: "OSMO",
    currVal: 4,
    prevVal: 4,
    changeVal: 0,
  },
  {
    chainName: "Polygon",
    tokenName: "MATIC",
    currVal: 2,
    prevVal: 2,
    changeVal: 0,
  },
  {
    chainName: "Mina",
    tokenName: "MINA",
    currVal: 11,
    prevVal: 11,
    changeVal: 0,
  },
];

export default function getData(...lens) {
    return data;
}
