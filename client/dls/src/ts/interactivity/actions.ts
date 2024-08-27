import { ctx, sceneManager, view } from '@src/routes/dls/+page.svelte';
import {
	type Circuit,
	cloneGraphAfterCircuit,
	CircuitSceneObject,
	getPropSetter
} from '@ts/scene/objects/circuits/circuit';
import { Wire, type WireAttachmentData } from '@ts/scene/objects/wire';
import type { UserAction } from './actions-manager';
import { Scene, type ID } from '@ts/scene/scene';
import { currentScene } from '@stores/currentScene';
import type { Vec2 } from '@ts/math';
import { circuitProps, focusedCircuit } from '@src/lib/stores/focusedCircuit';
import { domLog } from '@src/lib/stores/debugging';
import { icNames } from '@src/lib/stores/integrated-circuits';
import { icInstantiators, icInstanciator } from '@src/lib/stores/circuitInstantiators';
import { integratedCircuits } from '@src/lib/stores/integrated-circuits';
import type { View } from '@ts/view-manager';
import { actionURL } from '@ts/api/helpers';
import { combineHeaders, getAuthHeader, getProjectIDHeader } from '../helpers';

export const DUMMY_HOSTNAME = 'this-url-should-not-be-fetched';
export const DUMMY_URL = new URL(`http://${DUMMY_HOSTNAME}/`);

export const clipboard = {
	circuits: new Array<Circuit>(),
	wires: new Array<Wire>()
};

function convertWireForApi(id: ID, attachmentData: WireAttachmentData) {
	const fromSceneObject = attachmentData.from.parentCircuit.sceneObject;
	if (fromSceneObject == null || fromSceneObject.id == null) {
		throw Error();
	}
	const toSceneObject = attachmentData.to.parentCircuit.sceneObject;
	if (toSceneObject == null || toSceneObject.id == null) {
		throw Error();
	}

	return {
		id,
		fromCircuit: fromSceneObject.id,
		fromPin: attachmentData.from.pinIndex,
		toCircuit: toSceneObject.id,
		toPin: attachmentData.to.pinIndex
	};
}

function convertCircuitForApi(id: ID, sceneObject: CircuitSceneObject) {
	const circuit = sceneObject.parentCircuit;
	const inputWires = circuit.consumerPins
		.map((pin) => {
			if (pin.wire == null) {
				return undefined;
			}

			if (pin.wire.id == null) {
				throw Error();
			}

			return convertWireForApi(pin.wire.id, pin.wire.attachmentData());
		})
		.filter((val) => val != null);

	const outputWires = circuit.producerPins.flatMap((pin) => {
		return pin.wires.map((wire) => {
			if (wire.id == null) {
				throw Error();
			}

			return convertWireForApi(wire.id, wire.attachmentData());
		});
	});

	return {
		id,
		// 	ID          types.IDType       `json:"id"`

		circuitType: sceneObject.parentCircuit.circuitType,
		// 	CircuitType string             `binding:"required" json:"circuitType"`
		posWrl: sceneObject.posWrl.clone(),
		// 	PosWrl      math.Vec2          `binding:"required" json:"posWrl"`
		nInputPins: sceneObject.parentCircuit.consumerPins.length,
		// 	NInputPins  uint64             `json:"nInputPins"`
		nOutputPins: sceneObject.parentCircuit.producerPins.length,
		// 	NOutputPins uint64             `json:"nOutputPins"`

		inputWires,
		// 	InputWires  []*Wire            `json:"inputWires"`
		outputWires,
		// 	OutputWires []*Wire            `json:"outputWires"`
		props: Object.assign({}, circuit.props)
		// 	Props       types.CircuitProps `json:"props"`
	};
}

export function cloneSelectedAfterCircuit(
	start: Circuit,
	clonedCircuits: Circuit[],
	clonedWires: Wire[],
	circuitCloneMapping: Map<Circuit, Circuit>,
	wireCloneMapping: Map<Wire, Wire>
) {
	if (start.sceneObject == null || !start.sceneObject.isSelected) {
		return undefined;
	}

	const tmp = circuitCloneMapping.get(start);
	if (tmp != null) {
		return tmp;
	}

	const circuit = start;
	const cloned = circuit.clone();

	clonedCircuits.push(cloned);
	circuitCloneMapping.set(circuit, cloned);

	for (let pPinIdx = 0; pPinIdx < circuit.producerPins.length; pPinIdx++) {
		for (let wireIdx = 0; wireIdx < circuit.producerPins[pPinIdx].wires.length; wireIdx++) {
			console.log('[cloneCircuitTree] pPinIdx: ', pPinIdx);
			console.log('[cloneCircuitTree] circuit: ', circuit);
			console.log('[cloneCircuitTree] cloned: ', cloned);
			const clonedWire = cloneSelectedAfterWire(
				circuit.producerPins[pPinIdx].wires[wireIdx],
				clonedCircuits,
				clonedWires,
				circuitCloneMapping,
				wireCloneMapping
			);
			if (clonedWire == null) {
				cloned.producerPins[pPinIdx].wires.splice(wireIdx, 1);
			} else {
				cloned.producerPins[pPinIdx].wires[wireIdx] = clonedWire;
			}
		}
	}
	return cloned;
}

function cloneSelectedAfterWire(
	start: Wire,
	clonedCircuits: Circuit[],
	clonedWires: Wire[],
	circuitCloneMapping: Map<Circuit, Circuit>,
	wireCloneMapping: Map<Wire, Wire>
) {
	const producerPin = start.producerPin;
	const consumerPin = start.consumerPin;

	if (
		producerPin == null ||
		producerPin.parentCircuit.sceneObject == null ||
		!producerPin.parentCircuit.sceneObject.isSelected ||
		consumerPin == null ||
		consumerPin.parentCircuit.sceneObject == null ||
		!consumerPin.parentCircuit.sceneObject.isSelected
	) {
		return undefined;
	}

	const tmp = wireCloneMapping.get(start);
	if (tmp != null) {
		return tmp;
	}

	const wire = start;
	const cloned = wire.clone();

	clonedWires.push(cloned);
	wireCloneMapping.set(wire, cloned);

	const consumerCircuit = cloneGraphAfterCircuit(
		consumerPin.parentCircuit,
		clonedCircuits,
		clonedWires,
		circuitCloneMapping,
		wireCloneMapping
	);

	// console.log("[cloneCircuitTree] [Wire] id: ", start.id);
	// console.log("[cloneCircuitTree] [Wire] wire: ", wire);
	// console.log("[cloneCircuitTree] [Wire] cloned: ", cloned);
	cloned.setConsumerPinNoUpdate(consumerCircuit.consumerPins[consumerPin.pinIndex]);

	const producerCircuit = cloneGraphAfterCircuit(
		producerPin.parentCircuit,
		clonedCircuits,
		clonedWires,
		circuitCloneMapping,
		wireCloneMapping
	);
	cloned.setProducerPinNoUpdate(producerCircuit.producerPins[producerPin.pinIndex]);
	return cloned;
}

export function copySelectedToClipboard() {
	const clonedCircuits = new Array<Circuit>();
	const clonedWires = new Array<Wire>();

	const circuitCloneMapping = new Map<Circuit, Circuit>();
	const wireCloneMapping = new Map<Wire, Wire>();

	for (const circuit of sceneManager.selectedCircuits) {
		cloneSelectedAfterCircuit(
			circuit.parentCircuit,
			clonedCircuits,
			clonedWires,
			circuitCloneMapping,
			wireCloneMapping
		);
	}

	clipboard.circuits = clonedCircuits;
	clipboard.wires = clonedWires;
	console.log('Clipboard: ', clipboard);
}

export function dragSelection(delta: Vec2) {
	for (const circuit of sceneManager.selectedCircuits) {
		circuit.setPos(circuit.tightRectWrl.xy.add(delta));
	}
}

export class PasteFromClipBoardUserAction implements UserAction {
	name = 'PasteFromClipBoardUserAction';

	clonedCircuits: Array<Circuit> = new Array();
	clonedWires: Array<Wire> = new Array();

	circuitIDs: Map<ID, CircuitSceneObject> | undefined = undefined;
	wireIDs: Map<ID, Wire> | undefined = undefined;

	targetSceneID: number;

	constructor(circuits: Array<Circuit>, wires: Array<Wire>) {
		const circuitCloneMapping = new Map<Circuit, Circuit>();
		const wireCloneMapping = new Map<Wire, Wire>();
		for (const circuit of circuits) {
			cloneGraphAfterCircuit(
				circuit,
				this.clonedCircuits,
				this.clonedWires,
				circuitCloneMapping,
				wireCloneMapping
			);
		}
		const currentScene_ = currentScene.get();
		if (currentScene_ == null || currentScene_.id == null) {
			throw Error();
		}
		this.targetSceneID = currentScene_.id;
	}

	do(): void {
		sceneManager.deselectAll();

		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}

		if (this.wireIDs == null) {
			this.wireIDs = new Map();
			for (const wire of this.clonedWires) {
				const id = currentScene_.getNextID();
				currentScene_.registerWireWithId(id, wire);
				this.wireIDs.set(id, wire);
			}
		} else {
			for (const [wireID, wire] of this.wireIDs.entries()) {
				currentScene_.registerWireWithId(wireID, wire);
			}
		}

		if (this.circuitIDs == null) {
			this.circuitIDs = new Map();
			for (const circuit of this.clonedCircuits) {
				if (circuit.sceneObject == null) {
					throw Error();
				}

				const id = currentScene_.getNextID();
				const sceneObject = CircuitSceneObject.newRegistered(
					id,
					circuit,
					circuit.sceneObject.tightRectWrl.xy,
					currentScene_,
					ctx
				);
				this.circuitIDs.set(id, sceneObject);
				sceneManager.selectCircuit(circuit.sceneObject.id as ID);
			}
		} else {
			for (const [id, circuit] of this.circuitIDs.entries()) {
				currentScene_.registerCircuitWithID(id, circuit);
				sceneManager.selectCircuit(id);
			}
		}
	}
	undo(): void {
		if (this.circuitIDs == null || this.wireIDs == null) {
			throw Error();
		}
		const currentScene = sceneManager.getCurrentScene();
		if (currentScene == null) {
			throw Error();
		}
		for (const [id] of this.circuitIDs.entries()) {
			currentScene.unregisterCircuit(id);
		}
		for (const [id] of this.wireIDs.entries()) {
			currentScene.unregisterWire(id);
		}
		sceneManager.deselectAll();
	}
	async hitDoEndpoint(): Promise<Response> {
		if (this.circuitIDs == null || this.wireIDs == null) {
			throw Error();
		}
		const circuits = Array.from(this.circuitIDs.entries()).map(([id, circuit]) => {
			const apicircuit = convertCircuitForApi(id, circuit);
			console.log('Original circuit', circuit);
			console.log('Api circuit: ', apicircuit);
			return apicircuit;
		});
		const wires = Array.from(this.wireIDs.entries()).map(([id, wire]) =>
			convertWireForApi(id, wire.attachmentData())
		);
		return await customFetch(actionURL('paste-from-clipboard/do'), {
			body: JSON.stringify({
				targetSceneID: this.targetSceneID,
				circuits,
				wires
			})
		});
	}
	async hitUndoEndpoint(): Promise<Response> {
		if (this.circuitIDs == null || this.wireIDs == null) {
			throw Error();
		}
		const wireIDs = Array.from(this.wireIDs.keys());
		const circuitIDs = Array.from(this.circuitIDs.keys());
		return await customFetch(actionURL('paste-from-clipboard/undo'), {
			body: JSON.stringify({
				targetSceneID: this.targetSceneID,
				circuitIDs,
				wireIDs
			})
		});
	}
}

async function customFetch(url: URL, init: { body: any }) {
	return await fetch(url, {
		method: 'POST',
		body: init.body,
		credentials: 'include',
		headers: combineHeaders(getAuthHeader(), getProjectIDHeader())
	});
}

// Api implemented
export class NoopUserAction implements UserAction {
	name = 'Noop';
	do(): void {}
	undo(): void {}

	async hitDoEndpoint() {
		return await customFetch(actionURL('noop'), { body: JSON.stringify(this) });
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('noop'), { body: JSON.stringify(this) });
	}
}

// Api Implemented
export class DragUserAction implements UserAction {
	name = 'Drag';
	constructor(private deltaWrl: Vec2) {
		// console.log('Drag Action Created');
	}

	do(): void {
		dragSelection(this.deltaWrl);
	}
	undo(): void {
		dragSelection(this.deltaWrl.neg());
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('/drag-selection/do'), {
			body: JSON.stringify(this)
		});
	}

	async hitUndoEndpoint() {
		return await customFetch(actionURL('/drag-selection/undo'), {
			body: JSON.stringify(this)
		});
	}
}

// Api implemented
export class PanUserAction implements UserAction {
	constructor(private deltaScr: Vec2) {}
	name = 'Pan';

	do(): void {
		view.pan(this.deltaScr);
	}
	undo(): void {
		view.pan(this.deltaScr.neg());
	}
	async hitDoEndpoint() {
		// return DUMMY_URL;
		return await customFetch(actionURL('/pan/do'), { body: JSON.stringify(this) });
	}
	async hitUndoEndpoint() {
		// return DUMMY_URL;
		return await customFetch(actionURL('/pan/undo'), { body: JSON.stringify(this) });
	}
}

// Api implemented
export class ZoomUserAction implements UserAction {
	constructor(
		readonly zoomOriginScr: Vec2,
		public zoomLevelDelta: number
	) {}

	name = 'Zoom';

	do(): void {
		view.zoom(this.zoomOriginScr, view.zoomLevel + this.zoomLevelDelta);
	}
	undo(): void {
		view.zoom(this.zoomOriginScr, view.zoomLevel - this.zoomLevelDelta);
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('/mouse-zoom/do'), { body: JSON.stringify(this) });
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('/mouse-zoom/undo'), {
			body: JSON.stringify(this)
		});
	}
}

// Api implemented
export class TouchScreenZoomUserAction implements UserAction {
	name = '';
	constructor(
		readonly startingView: View,
		readonly endingView: View
	) {}

	do(): void {
		view.setView(this.endingView);
	}
	undo(): void {
		view.setView(this.startingView);
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('/touch-screen-zoom/do'), {
			body: JSON.stringify(this)
		});
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('/touch-screen-zoom/undo'), {
			body: JSON.stringify(this)
		});
	}
}

// Api implemented
export class DeleteSelectionUserAction implements UserAction {
	name = 'DeleteSelectionUserAction';

	private circuits: Map<ID, CircuitSceneObject> = new Map();
	private selectedWires: Map<ID, Wire> = new Map();
	private selectedWiresAttachmentData: Map<ID, WireAttachmentData> = new Map();
	private unSelectedWires: Map<ID, Wire> = new Map();
	private unselectedWiresAttachmentData: Map<ID, WireAttachmentData> = new Map();

	constructor(private sceneID: ID) {
		const scene = sceneManager.scenes.get(sceneID);
		if (scene == null) {
			throw Error();
		}

		for (const circuit of sceneManager.selectedCircuits) {
			if (circuit.id == null) {
				throw Error();
			}
			const parentCircuit = circuit.parentCircuit;

			for (const cPin of parentCircuit.consumerPins) {
				if (cPin.wire != null) {
					if (cPin.wire.id == null) {
						throw Error();
					}

					if (!cPin.wire.isSelected && !this.unSelectedWires.has(cPin.wire.id)) {
						this.unSelectedWires.set(cPin.wire.id, cPin.wire);
						if (cPin.wire.producerPin == null) {
							throw Error();
						}
						if (cPin.wire.consumerPin == null) {
							throw Error();
						}

						this.unselectedWiresAttachmentData.set(cPin.wire.id, {
							from: cPin.wire.producerPin,
							to: cPin.wire.consumerPin
						});
					} else if (cPin.wire.isSelected && !this.selectedWires.has(cPin.wire.id)) {
						this.selectedWires.set(cPin.wire.id, cPin.wire);
						if (cPin.wire.producerPin == null) {
							throw Error();
						}
						if (cPin.wire.consumerPin == null) {
							throw Error();
						}

						this.selectedWiresAttachmentData.set(cPin.wire.id, {
							from: cPin.wire.producerPin,
							to: cPin.wire.consumerPin
						});
					}
				}
			}

			for (const pPin of parentCircuit.producerPins) {
				for (const wire of pPin.wires) {
					if (wire.id == null) {
						throw Error();
					}
					if (!wire.isSelected && !this.unSelectedWires.has(wire.id)) {
						this.unSelectedWires.set(wire.id, wire);
						if (wire.producerPin == null) {
							throw Error();
						}
						if (wire.consumerPin == null) {
							throw Error();
						}
						this.unselectedWiresAttachmentData.set(wire.id, {
							from: wire.producerPin,
							to: wire.consumerPin
						});
					} else if (wire.isSelected && !this.selectedWires.has(wire.id)) {
						this.selectedWires.set(wire.id, wire);
						if (wire.producerPin == null) {
							throw Error();
						}
						if (wire.consumerPin == null) {
							throw Error();
						}
						this.selectedWiresAttachmentData.set(wire.id, {
							from: wire.producerPin,
							to: wire.consumerPin
						});
					}
				}
			}

			this.circuits.set(circuit.id, circuit);
		}

		for (const wire of sceneManager.selectedWires) {
			if (wire.id == null) {
				throw Error();
			}
			if (!this.selectedWires.has(wire.id)) {
				this.selectedWires.set(wire.id, wire);
				if (wire.consumerPin == null) {
					throw Error();
				}
				if (wire.producerPin == null) {
					throw Error();
				}

				this.selectedWiresAttachmentData.set(wire.id, {
					from: wire.producerPin,
					to: wire.consumerPin
				});
			}
		}
	}

	do(): void {
		const scene = sceneManager.scenes.get(this.sceneID);
		if (scene == null) {
			throw Error();
		}
		for (const [id, circuit] of this.circuits.entries()) {
			scene.unregisterCircuit(id);
			sceneManager.selectedCircuits.delete(circuit);
		}
		for (const wire of this.selectedWires.values()) {
			wire.detach();
			sceneManager.selectedWires.delete(wire);
		}
		for (const wire of this.unSelectedWires.values()) {
			wire.detach();
		}
	}
	undo(): void {
		const scene = sceneManager.scenes.get(this.sceneID);
		if (scene == null) {
			throw Error();
		}
		for (const [id, circuit] of this.circuits) {
			scene.registerCircuitWithID(id, circuit);
			if (circuit.isSelected) {
				sceneManager.selectCircuitUnchecked(id);
			}
		}
		for (const [id, wire] of this.selectedWires) {
			const attachmentData = this.selectedWiresAttachmentData.get(id);
			if (attachmentData == null) {
				throw Error();
			}
			scene.registerWireWithId(id, wire);
			wire.attach(attachmentData);
			sceneManager.selectWireUnchecked(wire);
		}
		for (const [id, wire] of this.unSelectedWires) {
			const attachmentData = this.unselectedWiresAttachmentData.get(id);
			if (attachmentData == null) {
				throw Error();
			}
			scene.registerWireWithId(id, wire);
			wire.attach(attachmentData);
		}
	}

	async hitDoEndpoint(): Promise<Response> {
		const circuitIDs: ID[] = [];
		for (const id of this.circuits.keys()) {
			circuitIDs.push(id);
		}

		const wireIDs: ID[] = [];
		for (const id of this.selectedWires.keys()) {
			wireIDs.push(id);
		}
		for (const id of this.unSelectedWires.keys()) {
			wireIDs.push(id);
		}

		return await customFetch(actionURL('delete-selected/do'), {
			body: JSON.stringify({
				sceneID: this.sceneID,
				circuitIDs,
				wireIDs
			})
		});
	}
	async hitUndoEndpoint(): Promise<Response> {
		const circuits = Array.from(this.circuits.entries()).map(([id, sceneObject]) => {
			return convertCircuitForApi(id, sceneObject);
		});
		const selectedWires = Array.from(this.selectedWires.entries()).map(([id, wire]) => {
			return convertWireForApi(id, wire.attachmentData());
		});
		const unSelectedWires = Array.from(this.unSelectedWires.entries()).map(([id, wire]) => {
			return convertWireForApi(id, wire.attachmentData());
		});
		return await customFetch(actionURL('delete-selected/undo'), {
			body: JSON.stringify({
				sceneID: this.sceneID,
				circuits,
				selectedWires,
				unSelectedWires
			})
		});
	}
}

// Api implemented
export class CreateCircuitUserAction implements UserAction {
	name = 'CreateCircuitUserAction';
	private readonly circuitID: ID;
	constructor(
		private targetSceneID: ID,

		// do NOT remove this property, this is a parameter for the api.
		// It will not be used on the client side.
		private circuitType: string,

		private instantiator: () => Circuit,
		private locScr: Vec2
	) {
		const targetScene = sceneManager.scenes.get(this.targetSceneID);
		if (targetScene == null) {
			throw Error();
		}
		this.circuitID = targetScene.getNextID();
	}

	do(): void {
		console.log('CreateCircuitUserAction.do');
		const targetScene = sceneManager.scenes.get(this.targetSceneID);
		if (targetScene == null) {
			throw Error();
		}
		const circuit = this.instantiator();
		const currentScene = sceneManager.getCurrentScene();

		CircuitSceneObject.newRegistered(
			this.circuitID,
			circuit,
			view.screenToWorld(this.locScr),
			currentScene,
			ctx
		);
	}
	undo(): void {
		console.log('CreateCircuitUserAction.undo');
		const targetScene = sceneManager.scenes.get(this.targetSceneID);
		if (targetScene == null) {
			throw Error();
		}
		targetScene.unregisterCircuit(this.circuitID);
		console.log('TargetScene: ', targetScene);
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('/create-circuit/do'), {
			body: JSON.stringify(this)
		});
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('/create-circuit/undo'), {
			body: JSON.stringify(this)
		});
	}
}

// Api Implemented
export class CreateWireUserAction implements UserAction {
	name = 'CreateWireUserAction';
	wireID: ID;

	producerCircuitID: ID;
	producerPinIdx: number;

	consumerCircuitID: ID;
	consumerPinIdx: number;

	constructor(
		private sceneID: ID,
		wire: Wire
	) {
		const targetScene = sceneManager.scenes.get(sceneID);
		if (targetScene == null) {
			throw Error();
		}
		this.wireID = targetScene.getNextID();
		if (
			wire.consumerPin == null ||
			wire.consumerPin.parentCircuit.sceneObject == null ||
			wire.consumerPin.parentCircuit.sceneObject.id == null ||
			wire.producerPin == null ||
			wire.producerPin.parentCircuit.sceneObject == null ||
			wire.producerPin.parentCircuit.sceneObject.id == null
		) {
			throw Error();
		}

		this.producerCircuitID = wire.producerPin.parentCircuit.sceneObject.id;
		this.producerPinIdx = wire.producerPin.pinIndex;
		this.consumerCircuitID = wire.consumerPin.parentCircuit.sceneObject.id;
		this.consumerPinIdx = wire.consumerPin.pinIndex;
	}

	do(): void {
		const targetScene = sceneManager.scenes.get(this.sceneID);
		if (targetScene == null) {
			throw Error();
		}

		const producerCircuit = targetScene.idToCircuit.get(this.producerCircuitID);
		const consumerCircuit = targetScene.idToCircuit.get(this.consumerCircuitID);
		if (producerCircuit == null || consumerCircuit == null) {
			throw Error();
		}

		const wire = Wire.newUnregistered(
			producerCircuit.parentCircuit.producerPins[this.producerPinIdx],
			consumerCircuit.parentCircuit.consumerPins[this.consumerPinIdx]
		);

		wire.registerWithID(this.wireID, targetScene);
		console.log('[CreatingWireUserAction.do] wire: ', wire);
	}
	undo(): void {
		const targetScene = sceneManager.scenes.get(this.sceneID);
		if (targetScene == null) {
			throw Error();
		}
		const wire = targetScene.idToWire.get(this.wireID);
		if (wire == null) {
			throw Error();
		}
		targetScene.unregisterWire(this.wireID);
		wire.detach();
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('/create-wire/do'), {
			body: JSON.stringify(this)
		});
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('/create-wire/undo'), {
			body: JSON.stringify(this)
		});
	}
}

// Api implemented
export class DeleteWireUserAction implements UserAction {
	name = 'DeleteWireUserAction';

	wireID: ID;

	producerCircuitID: ID;
	producerPinIdx: number;

	consumerCircuitID: ID;
	consumerPinIdx: number;

	constructor(
		wire: Wire,
		private sceneID: ID
	) {
		const targetScene = sceneManager.scenes.get(sceneID);
		if (targetScene == null) {
			throw Error();
		}
		if (wire.id == null) {
			throw Error();
		}

		this.wireID = wire.id;

		if (
			wire.consumerPin == null ||
			wire.consumerPin.parentCircuit.sceneObject == null ||
			wire.consumerPin.parentCircuit.sceneObject.id == null ||
			wire.producerPin == null ||
			wire.producerPin.parentCircuit.sceneObject == null ||
			wire.producerPin.parentCircuit.sceneObject.id == null
		) {
			throw Error();
		}

		this.producerCircuitID = wire.producerPin.parentCircuit.sceneObject.id;
		this.producerPinIdx = wire.producerPin.pinIndex;
		this.consumerCircuitID = wire.consumerPin.parentCircuit.sceneObject.id;
		this.consumerPinIdx = wire.consumerPin.pinIndex;
	}
	do(): void {
		const targetScene = sceneManager.scenes.get(this.sceneID);
		if (targetScene == null) {
			throw Error();
		}
		const wire = targetScene.idToWire.get(this.wireID);
		if (wire == null) {
			throw Error();
		}
		targetScene.unregisterWire(this.wireID);
		wire.detach();
	}
	undo(): void {
		const targetScene = sceneManager.scenes.get(this.sceneID);
		if (targetScene == null) {
			throw Error();
		}

		const producerCircuit = targetScene.idToCircuit.get(this.producerCircuitID);
		const consumerCircuit = targetScene.idToCircuit.get(this.consumerCircuitID);
		if (producerCircuit == null || consumerCircuit == null) {
			throw Error();
		}

		const wire = Wire.newUnregistered(
			producerCircuit.parentCircuit.producerPins[this.producerPinIdx],
			consumerCircuit.parentCircuit.consumerPins[this.consumerPinIdx]
		);

		wire.registerWithID(this.wireID, targetScene);
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('/delete-wire/do'), {
			body: JSON.stringify(this)
		});
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('/delete-wire/undo'), {
			body: JSON.stringify(this)
		});
	}
}

export class SetCircuitPropUserAction implements UserAction {
	name = 'SetCircuitPropUserAction';

	constructor(
		private sceneID: ID,
		private circuitID: ID,
		private propName: string,
		private valueToSet: any,
		private currentValue: any
	) {
		const scene = sceneManager.scenes.get(sceneID);
		if (scene == null) {
			throw Error();
		}
		const circuit = scene.idToCircuit.get(circuitID);
		if (circuit == null) {
			throw Error();
		}
	}
	do(): void {
		const scene = sceneManager.scenes.get(this.sceneID);
		if (scene == null) {
			throw Error();
		}
		const circuit = scene.idToCircuit.get(this.circuitID);
		if (circuit == null) {
			throw Error();
		}

		const propSetter = getPropSetter(circuit.parentCircuit, this.propName);
		propSetter(circuit.parentCircuit, this.valueToSet);
		circuitProps.refresh();
	}

	undo(): void {
		const scene = sceneManager.scenes.get(this.sceneID);
		if (scene == null) {
			throw Error();
		}
		const circuit = scene.idToCircuit.get(this.circuitID);
		if (circuit == null) {
			throw Error();
		}

		const propSetter = getPropSetter(circuit.parentCircuit, this.propName);
		propSetter(circuit.parentCircuit, this.currentValue);

		circuitProps.refresh();
	}
	async hitDoEndpoint() {
		const body = {
			sceneID: this.sceneID,
			circuitID: this.circuitID,
			propName: this.propName,
			valueToSet: String(this.valueToSet),
			currentValue: String(this.currentValue)
		};

		return await customFetch(actionURL('/set-prop/do'), { body: JSON.stringify(body) });
	}
	async hitUndoEndpoint() {
		const body = {
			sceneID: this.sceneID,
			circuitID: this.circuitID,
			propName: this.propName,
			valueToSet: String(this.valueToSet),
			currentValue: String(this.currentValue)
		};
		return await customFetch(actionURL('/set-prop/undo'), { body: JSON.stringify(body) });
	}
}

// Api Implemented
export class SelectCircuitUserAction implements UserAction {
	name = 'SelectCircuitUserAction';
	private sceneID: ID;
	constructor(private circuitID: ID) {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}
		const sceneID = currentScene_.id;
		if (sceneID == null) {
			throw Error();
		}
		this.sceneID = sceneID;

		const scene = sceneManager.scenes.get(this.sceneID);
		if (scene == null) {
			throw Error();
		}
		if (scene.idToCircuit.get(this.circuitID) == null) {
			throw Error();
		}
	}
	do(): void {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}

		if (currentScene_.id != this.sceneID) {
			throw Error();
		}
		sceneManager.selectCircuit(this.circuitID);
	}
	undo(): void {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}

		if (currentScene_.id != this.sceneID) {
			throw Error();
		}
		sceneManager.deselectCircuit(this.circuitID);
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('/select-circuit/do'), {
			body: JSON.stringify(this)
		});
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('/select-circuit/undo'), {
			body: JSON.stringify(this)
		});
	}
}

// Api Implemented
export class DeselectCircuitUserAction implements UserAction {
	name = 'DeselectCircuitUserAction';
	private sceneID: ID;
	constructor(private circuitID: ID) {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}
		const sceneID = currentScene_.id;
		if (sceneID == null) {
			throw Error();
		}
		this.sceneID = sceneID;

		const scene = sceneManager.scenes.get(this.sceneID);
		if (scene == null) {
			throw Error();
		}
		if (scene.idToCircuit.get(this.circuitID) == null) {
			throw Error();
		}
	}
	do(): void {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}
		if (currentScene_.id != this.sceneID) {
			throw Error();
		}
		sceneManager.deselectCircuit(this.circuitID);
	}
	undo(): void {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}

		if (currentScene_.id != this.sceneID) {
			throw Error();
		}
		sceneManager.selectCircuit(this.circuitID);
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('/deselect-circuit/do'), {
			body: JSON.stringify(this)
		});
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('/deselect-circuit/undo'), {
			body: JSON.stringify(this)
		});
	}
}

// Api Implemented
export class SwitchSceneUserAction implements UserAction {
	name = 'SwitchSceneUserAction';
	private fromSceneID: ID;

	private selectedCircuits = new Array<ID>();
	private selectedWires = new Array<ID>();

	constructor(private toSceneID: ID) {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}

		if (currentScene_.id == null) {
			throw Error();
		}
		this.fromSceneID = currentScene_.id;
		for (const circuit of sceneManager.selectedCircuits) {
			if (circuit.id == null) {
				throw Error();
			}
			this.selectedCircuits.push(circuit.id as ID);
		}
		for (const wire of sceneManager.selectedWires) {
			if (wire.id == null) {
				console.error('Wire without id: ', wire);
				throw Error();
			}
			this.selectedWires.push(wire.id as ID);
		}
	}

	do(): void {
		const toScene = sceneManager.scenes.get(this.toSceneID);
		if (toScene == null) {
			domLog('[SceneManager.setCurrentScene] Provided sceneId is invalid.');
			throw Error();
		}

		console.log('SwitchSceneUserAction.do');
		toScene.refreshICLabels();

		currentScene.set(toScene);
		sceneManager.deselectAll();
	}

	undo(): void {
		const fromScene = sceneManager.scenes.get(this.fromSceneID);
		if (fromScene == null) {
			domLog('[SceneManager.setCurrentScene] Provided sceneId is invalid.');
			throw Error();
		}

		console.log('SwitchSceneUserAction.undo');
		fromScene.refreshICLabels();

		currentScene.set(fromScene);

		sceneManager.deselectAllCircuits();
		const scene = sceneManager.scenes.get(this.fromSceneID);
		if (scene == null) {
			throw Error();
		}
		for (const circuitID of this.selectedCircuits) {
			sceneManager.selectCircuitUnchecked(circuitID);
		}

		if (this.selectedCircuits.length != 1) {
			focusedCircuit.set(undefined);
		} else {
			const circuitID = this.selectedCircuits[0];
			const circuit = scene.idToCircuit.get(circuitID);
			if (circuit == null) {
				throw Error();
			}

			focusedCircuit.set(circuit);
		}

		sceneManager.selectedWires = new Set();
		for (const wireID of this.selectedWires) {
			const wire = scene.idToWire.get(wireID);
			if (wire == null) {
				throw Error();
			}
			sceneManager.selectWireUnchecked(wire);
			// sceneManager.selectedWires.add(wire);
		}
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('switch-scene/do'), {
			body: JSON.stringify(this)
		});
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('switch-scene/undo'), {
			body: JSON.stringify(this)
		});
	}
}

// Api implemented
export class CreateICUserAction implements UserAction {
	name = 'CreateICUserAction';

	static nextICNumber = 0;

	currentICNumber: number;

	private scene: Scene | undefined;
	private sceneID: number;
	constructor() {
		this.sceneID = sceneManager.getNextSceneID();
		this.currentICNumber = CreateICUserAction.nextICNumber;
		CreateICUserAction.nextICNumber += 1;
	}

	do(): void {
		this.scene = Scene.newWithIO();

		if (this.currentICNumber === 0) {
			this.scene.name = 'New Circuit';
		} else {
			this.scene.name = `New Circuit (${this.currentICNumber})`;
		}

		sceneManager.registerSceneWithID(this.sceneID, this.scene);
		icNames.add(this.scene.name.toLowerCase());

		icInstantiators.newInstantiator(this.sceneID, icInstanciator(this.sceneID));

		// integratedCircuits.newIC()
		integratedCircuits.update((circuits) => {
			if (this.scene == null) {
				throw Error();
			}

			circuits.set(this.sceneID, this.scene.name);
			return circuits;
		});
	}
	undo(): void {
		sceneManager.unregisterScene(this.sceneID);
		integratedCircuits.update((circuits) => {
			if (this.scene == null) {
				throw Error();
			}
			icNames.delete(this.scene.name.toLowerCase());

			circuits.delete(this.sceneID);
			return circuits;
		});

		icInstantiators.removeInstantiator(this.sceneID);
	}
	async hitDoEndpoint() {
		if (this.scene == null) {
			throw Error();
		}

		const body = {
			sceneID: this.sceneID,
			icName: this.scene.name
		};

		return await customFetch(actionURL('/create-ic/do'), { body: JSON.stringify(body) });
	}
	async hitUndoEndpoint() {
		if (this.scene == null) {
			throw Error();
		}

		const body = {
			sceneID: this.sceneID,
			icName: this.scene.name
		};

		return await customFetch(actionURL('create-ic/undo'), { body: JSON.stringify(body) });
	}
}

// Api implemented
export class RenameICUserAction implements UserAction {
	name = 'RenameICUserAction';

	from: string;
	constructor(
		private id: ID,
		private to: string
	) {
		const scene = sceneManager.scenes.get(id);
		if (scene == null) {
			throw Error();
		}
		this.from = scene.name;
	}

	do(): void {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}

		integratedCircuits.rename(this.id, this.to);
		currentScene_.refreshICLabels();
	}
	undo(): void {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}

		integratedCircuits.rename(this.id, this.from);
		currentScene_.refreshICLabels();
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('/rename-ic/do'), { body: JSON.stringify(this) });
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('/rename-ic/undo'), {
			body: JSON.stringify(this)
		});
	}
}

// Api implemented
export class DeselectAllUserAction implements UserAction {
	name = 'DeselectAllUserAction';

	private currentSceneID: ID;
	private selectedWireIDs = new Array<ID>();
	private selectedCircuitIDs = new Array<ID>();

	constructor() {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}
		const currentSceneID = currentScene_.id;
		if (currentSceneID == null) {
			throw Error();
		}
		this.currentSceneID = currentSceneID;

		for (const circuit of sceneManager.selectedCircuits) {
			if (circuit.id == null) {
				throw Error();
			}
			this.selectedCircuitIDs.push(circuit.id);
		}
		for (const wire of sceneManager.selectedWires) {
			if (wire.id == null) {
				throw Error();
			}
			this.selectedWireIDs.push(wire.id);
		}
	}
	do(): void {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}
		if (currentScene_.id != this.currentSceneID) {
			throw Error();
		}
		sceneManager.deselectAll();
	}
	undo(): void {
		const currentScene_ = currentScene.get();
		if (currentScene_ == null) {
			throw Error();
		}

		if (currentScene_.id != this.currentSceneID) {
			throw Error();
		}

		for (const circuit of this.selectedCircuitIDs) {
			sceneManager.selectCircuitUnchecked(circuit);
		}
		for (const wireID of this.selectedWireIDs) {
			const wire = currentScene_.idToWire.get(wireID);
			if (wire == null) {
				throw Error();
			}

			sceneManager.selectWireUnchecked(wire);
		}
	}
	async hitDoEndpoint() {
		return await customFetch(actionURL('deselect-all/do'), {
			body: JSON.stringify(this)
		});
	}
	async hitUndoEndpoint() {
		return await customFetch(actionURL('deselect-all/undo'), {
			body: JSON.stringify(this)
		});
	}
}
