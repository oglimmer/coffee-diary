import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import StarRating from './StarRating.vue'

const FILLED = '★' // ★
const EMPTY = '☆' // ☆

describe('StarRating', () => {
  it('always renders five star buttons', () => {
    const wrapper = mount(StarRating, { props: { modelValue: 0 } })
    expect(wrapper.findAll('.star-btn')).toHaveLength(5)
  })

  it('fills stars up to the current value', () => {
    const wrapper = mount(StarRating, { props: { modelValue: 3 } })
    const stars = wrapper.findAll('.star-btn')
    expect(stars.map((s) => s.text())).toEqual([FILLED, FILLED, FILLED, EMPTY, EMPTY])
    expect(wrapper.findAll('.star-btn.filled')).toHaveLength(3)
  })

  it('shows the label matching the rating', () => {
    const wrapper = mount(StarRating, { props: { modelValue: 4 } })
    expect(wrapper.find('.rating-label').text()).toBe('Excellent')
  })

  it('renders no label for an out-of-range rating', () => {
    const wrapper = mount(StarRating, { props: { modelValue: 0 } })
    expect(wrapper.find('.rating-label').exists()).toBe(false)
  })

  it('emits update:modelValue with the clicked star value', async () => {
    const wrapper = mount(StarRating, { props: { modelValue: 2 } })
    await wrapper.findAll('.star-btn')[4].trigger('click') // 5th star
    expect(wrapper.emitted('update:modelValue')).toEqual([[5]])
  })

  it('does not emit or show a label when readonly', async () => {
    const wrapper = mount(StarRating, { props: { modelValue: 3, readonly: true } })
    const stars = wrapper.findAll('.star-btn')
    expect(stars[0].attributes('disabled')).toBeDefined()
    await stars[0].trigger('click')
    expect(wrapper.emitted('update:modelValue')).toBeUndefined()
    expect(wrapper.find('.rating-label').exists()).toBe(false)
  })

  it('applies the size prop to button font-size', () => {
    const wrapper = mount(StarRating, { props: { modelValue: 1, size: 30 } })
    expect(wrapper.find('.star-btn').attributes('style')).toContain('font-size: 30px')
  })
})
