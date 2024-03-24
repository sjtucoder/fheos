import { AddCaller, AddCallee } from "../types/tests/contracts/Tx.sol";
import { createFheInstance, deployContract } from "./utils";
import {ethers} from "hardhat";

describe("Test Transactions Scenarios", () => {
  let contractCaller: AddCaller;
  let contractCallee: AddCallee;
  let contractAddr: string;

  beforeAll(async () => {
    contractCallee = (await deployContract("AddCallee")) as AddCallee;
    expect(contractCallee).toBeTruthy();

    contractCaller = (await deployContract("AddCaller", [contractCallee])) as AddCaller;
    expect(contractCaller).toBeTruthy();
    contractAddr = await contractCaller.getAddress();
    expect(contractAddr).toBeTruthy();
  });

  beforeEach(async () => {
    const encCounterResponse = await contractCaller.resetCounter();
    await encCounterResponse.wait();
  });

  it("Contract Deployment", async () => {
    contractCallee = (await deployContract("AddCallee")) as AddCallee;
    expect(contractCallee).toBeTruthy();

    contractCaller = (await deployContract("AddCaller", [contractCallee])) as AddCaller;
    expect(contractCaller).toBeTruthy();
    contractAddr = await contractCaller.getAddress();
    expect(contractAddr).toBeTruthy();
  });

  it("Basic Add Tx", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(8);

    // 1 - static call
    const encCounter = await contractCaller.addTx.staticCall(encInput, permit.publicKey);
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(8);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.addTx(encInput, permit.publicKey);
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(permit.publicKey);
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(8);
  });

  // todo: this fails!
  it.skip("Add via contract call as Plaintext", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);

    // 1 - static call
    const encCounter = await contractCaller.addViaContractCallAsPlain.staticCall(9, permit.publicKey);
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(9);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.addViaContractCallAsPlain(9, permit.publicKey);
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(permit.publicKey);
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(9);
  });

  it.only("Add via contract call as Encrypted Input", async () => {
    const { instance, permit } = await createFheInstance(contractAddr);
    const encInput = await instance.encrypt_uint32(10);

    // 1 - static call
    const encCounter = await contractCaller.addViaContractCall.staticCall(encInput, permit.publicKey);
    const counterStatic = instance.unseal(contractAddr, encCounter);
    expect(Number(counterStatic)).toEqual(10);

    // 2 - real call + query
    const encCounterReceipt = await contractCaller.addViaContractCall(encInput, permit.publicKey);
    await encCounterReceipt.wait();

    const getCounterResponse = await contractCaller.getCounter(permit.publicKey);
    const counter = instance.unseal(contractAddr, getCounterResponse);
    expect(Number(counter)).toEqual(10);
  });
});
